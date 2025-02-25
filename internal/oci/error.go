package oci

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

// APIErrorCode holds known / well-defined OCI API errors.
type APIErrorCode string

const (
	// APIErrorCodeBlobUnknown is returned when a blob is unknown to the
	// registry.
	APIErrorCodeBlobUnknown = "BLOB_UNKNOWN"
	// APIErrorCodeBlobUploadInvalid is returned when a blob upload is
	// invalid.
	APIErrorCodeBlobUploadInvalid = "BLOB_UPLOAD_INVALID"
	// APIErrorCodeBlobUploadUnknown is returned when a blob upload is
	// unknown to the registry.
	APIErrorCodeBlobUploadUnknown = "BLOB_UPLOAD_UNKNOWN"
	// APIErrorCodeDigestInvalid is returned when a provided digest did
	// not match uploaded content.
	APIErrorCodeDigestInvalid = "DIGEST_INVALID"
	// APIErrorCodeManifestBlobUnknwon is returned when a manifest
	// references a manifest or blob that is unknown to the registry.
	APIErrorCodeManifestBlobUnknwon = "MANIFEST_BLOB_UNKNOWN"
	// APIErrorCodeManifestInvalid is returned when a manifest is
	// invalid.
	APIErrorCodeManifestInvalid = "MANIFEST_INVALID"
	// APIErrorCodeManifestUnknown is returned when a manifest is unknown
	// to the registry.
	APIErrorCodeManifestUnknown = "MANIFEST_UNKNOWN"
	// APIErrorCodeNameInvalid is returned when an invalid repository
	// name is used.
	APIErrorCodeNameInvalid = "NAME_INVALID"
	// APIErrorCodeNameUnknown is returned when a repository name is not
	// known to registry.
	APIErrorCodeNameUnknown = "NAME_UNKNOWN"
	// APIErrorCodeSizeInvalid is returned when a provided length did not
	// match content length.
	APIErrorCodeSizeInvalid = "SIZE_INVALID"
	// APIErrorCodeUnauthorized is returned when authentication required.
	APIErrorCodeUnauthorized = "UNAUTHORIZED"
	// APIErrorCodeDenied is returned when the requested access to the
	// resource is denied.
	APIErrorCodeDenied = "DENIED"
	// APIErrorCodeUnsupported is returned when the operation is
	// unsupported.
	APIErrorCodeUnsupported = "UNSUPPORTED"
	// APIErrorCodeTooManyRequests is returned when the client has sent
	// too many requests.
	APIErrorCodeTooManyRequests = "TOOMANYREQUESTS"
)

var _ error = (*APIError)(nil)

// APIError is a common error type returned by [Client].
type APIError struct {
	Status     string
	StatusCode int
	Code       APIErrorCode
	Message    string
	Detail     string
}

// Error implements error.
func (d APIError) Error() string {
	return fmt.Sprintf("oci: %s - %s", d.Code, d.Message)
}

// assertStatusCode returns an error if the response does not match the given
// status code. If possible, extra detail is extracted and one or more
// [APIError] is returned.
func assertStatusCode(r *http.Response, statusCode int) error {
	if r.StatusCode == statusCode {
		return nil
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return httputil.Error{
			Status:     r.Status,
			StatusCode: r.StatusCode,
		}
	}

	var response struct {
		Errors []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Detail  string `json:"detail"`
		}
	}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return httputil.Error{
			Status:     r.Status,
			StatusCode: r.StatusCode,
		}
	}

	if len(response.Errors) == 0 {
		return httputil.Error{
			Status:     r.Status,
			StatusCode: r.StatusCode,
		}
	}

	errs := make([]error, len(response.Errors))
	for i, err := range response.Errors {
		errs[i] = APIError{
			Status:     r.Status,
			StatusCode: r.StatusCode,
			Code:       APIErrorCode(err.Code),
			Message:    err.Message,
			Detail:     err.Detail,
		}
	}

	return errors.Join(errs...)
}

// ErrorIsResourceUnknown returns true if the error is an APIError which points
// at the error being that the resource (blob, manifest, name) us unknown to the
// registry.
func ErrorIsResourceUnknown(err error) bool {
	if err == nil {
		return false
	}

	var apiErr APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case APIErrorCodeBlobUnknown, APIErrorCodeManifestUnknown, APIErrorCodeNameUnknown:
			return true
		}
	}

	return false
}
