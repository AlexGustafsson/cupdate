package httputil

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type spanKeyType struct{}

var spanKey = spanKeyType{}

// InstrumentHandler instruments a http.Handler.
// Instrumented handlers are required to set the http.route attribute to the
// span if available.
//
// If CORS is in use, it's the wrapped handler's responsibility to allow the
// traceresponse header.
//
// SEE: https://opentelemetry.io/docs/specs/semconv/http/http-spans/#http-server-semantic-conventions
func InstrumentHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := instrumentRequest(r)
		defer span.End()

		// SEE: https://github.com/w3c/trace-context/blob/main/spec/21-http_response_header_format.md
		w.Header().Set("traceresponse", fmt.Sprintf("00-%s-%s-01", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()))

		statusRecorder := StatusRecorder{Writer: w}

		handler.ServeHTTP(&statusRecorder, r.Clone(ctx))

		statusCode := statusRecorder.StatusCode()
		if statusCode == 0 {
			// Implicit status
			statusCode = http.StatusOK
		}

		span.SetAttributes(semconv.HTTPResponseStatusCode(statusCode))

		// SEE: https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
		if statusCode >= 500 {
			span.SetStatus(codes.Error, "")
		}
	})
}

// SpanFromRequest returns the current span for a handler instrumented using
// [InstrumentHandler]. Panics if handler is not instrumented.
// A span received this way MUST NOT be manually closed. It is closed once the
// instrumented handler completed the request.
func SpanFromRequest(r *http.Request) (context.Context, trace.Span) {
	spanValue := r.Context().Value(spanKey)
	if spanValue != nil {
		span, ok := spanValue.(trace.Span)
		if ok {
			return r.Context(), span
		}
	}

	panic("httputil: request not instrumented")
}

func instrumentRequest(r *http.Request) (context.Context, trace.Span) {
	requestURL, err := ResolveRequestURL(r)
	if err != nil {
		requestURL = r.URL
	}

	serverPort := 0
	switch requestURL.Port() {
	case "":
		switch requestURL.Scheme {
		case "http":
			serverPort = 80
		case "https":
			serverPort = 443
		default:
			panic("httputil: client misue - unsupported protocol")
		}
	default:
		// Port should already be parsed
		v, _ := strconv.ParseInt(requestURL.Port(), 10, 16)
		serverPort = int(v)
	}

	ctx, span := otel.Tracer(otelutil.DefaultScope).Start(
		r.Context(),
		r.Method+" "+requestURL.Path,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPRoute(requestURL.Path),
			semconv.ServerAddress(requestURL.Hostname()),
			semconv.ServerPort(serverPort),
			semconv.URLPath(requestURL.Path),
			semconv.URLScheme(requestURL.Scheme),
		),
	)

	if requestURL.RawQuery != "" {
		span.SetAttributes(semconv.URLQuery(requestURL.RawQuery))
	}

	if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
		span.SetAttributes(semconv.UserAgentOriginal(userAgent))
	}

	host, portString, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		port, _ := strconv.ParseInt(portString, 10, 16)
		span.SetAttributes(
			// TODO: Resolve these using remote peer headers. It opens up for spoofing
			// attacks, but for the sake of traces/metrics, it shouldn't matter? It's
			// not for audit purposes...
			semconv.ClientAddress(host),
			semconv.ClientPort(int(port)),
			semconv.ClientSocketAddress(host),
			semconv.ClientSocketPort(int(port)),
		)
	}

	return context.WithValue(ctx, spanKey, span), span
}
