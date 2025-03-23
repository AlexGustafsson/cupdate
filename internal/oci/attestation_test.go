package oci

import (
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttestationUnmarshalJSON(t *testing.T) {
	fs, err := os.OpenRoot("./testdata/attestations")
	require.NoError(t, err)

	testCases := []struct {
		// Path is the path within fs
		Path     string
		Expected *Attestation
	}{
		{
			Path: "1.json",
			Expected: &Attestation{
				BuildStartedOn:  time.Date(2025, 03, 22, 11, 23, 15, 399801890, time.UTC),
				BuildFinishedOn: time.Date(2025, 03, 22, 11, 26, 10, 631712748, time.UTC),
				Source:          "https://github.com/AlexGustafsson/cupdate",
				SourceRevision:  "2fbefbc45dd73e49c981e7a59c9b3d65314ed315",
				Dockerfile: `FROM --platform=${BUILDPLATFORM} node:22.14.0@sha256:cfef4432ab2901fd6ab2cb05b177d3c6f8a7f48cb22ad9d7ae28bb6aa5f8b471 AS web-builder

WORKDIR /src

COPY .yarnrc.yml package.json yarn.lock .
COPY .yarn .yarn

RUN yarn install

COPY tsconfig.json vite.config.ts .
COPY web web

ARG CUPDATE_VERSION="development build"
RUN VITE_CUPDATE_VERSION="${CUPDATE_VERSION}" yarn build

FROM --platform=${BUILDPLATFORM} golang:1.24.1@sha256:c5adecdb7b3f8c5ca3c88648a861882849cc8b02fed68ece31e25de88ad13418 AS builder

WORKDIR /src

# Use the toolchain specified in go.mod, or newer
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum .
RUN go mod download && go mod verify

COPY cmd cmd
COPY internal internal

COPY --from=web-builder /src/internal/web/public /src/internal/web/public

ARG CUPDATE_VERSION="development build"
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=${TARGETARCH} GOOS=${TARGETOS} CGO_ENABLED=0 go build -a -ldflags="-s -w -X 'main.Version=$CUPDATE_VERSION'" -o cupdate cmd/cupdate/*.go

FROM scratch AS export

COPY --from=builder /src/cupdate cupdate

FROM export

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV PATH=/

ENTRYPOINT ["cupdate"]
`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Path, func(t *testing.T) {
			file, err := fs.Open(testCase.Path)
			require.NoError(t, err)
			defer file.Close()

			content, err := io.ReadAll(file)
			require.NoError(t, err)

			var attestation Attestation
			require.NoError(t, json.Unmarshal(content, &attestation))

			assert.Equal(t, testCase.Expected, &attestation)
		})
	}
}
