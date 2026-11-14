FROM --platform=${BUILDPLATFORM} node:24.12.0@sha256:929c026d5a4e4a59685b3c1dbc1a8c3eb090aa95373d3a4fd668daa2493c8331 AS web-builder

RUN apt-get update -y && \
    apt-get install --no-install-recommends -y jq && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY .npmrc package.json package-lock.json .

RUN --mount=type=cache,target=node_modules  \
  npm ci

COPY tsconfig.json vite.config.ts .
COPY web web

ARG CUPDATE_VERSION="development build"
RUN --mount=type=cache,target=node_modules \
  VITE_CUPDATE_VERSION="${CUPDATE_VERSION}" npm run build && \
  npm sbom --omit dev --sbom-format spdx > sbom.spdx.json

# Properly format the frontend's root entry to make it clear that it's not an
# NPM package, but just a component used by the server
RUN cat << 'EOF' | /bin/sh
  jq '
  (.packages[0].SPDXID) as $id
  | .packages |= map(
      if .SPDXID == $id then
        .primaryPackagePurpose = "OTHER"
        | del(.versionInfo)
        | if has("externalRefs") then
            .externalRefs |= map(select(.referenceType != "purl"))
            | if (.externalRefs | length) == 0
              then del(.externalRefs)
              else .
              end
          else .
          end
      else .
      end
    )
  ' sbom.spdx.json > sbom2.spdx.json
  mv sbom2.spdx.json sbom.spdx.json
EOF

FROM --platform=${BUILDPLATFORM} golang:1.27.3@sha256:2920bb124f1a2e12a1d39f0cab756ae0068441ee1da7655dec7f5e374944095b AS osv-scanner-builder

ARG TARGETARCH
ARG TARGETOS

ARG OSV_SCANNER_REPO="https://github.com/google/osv-scanner"
# v2.4.0 (2026-06-18)
ARG OSV_SCANNER_REF="b56b5191101d5f27d4787d5583d8d01e9518a7af"

WORKDIR /src

# Use the toolchain specified in go.mod, or newer
ENV GOTOOLCHAIN=auto

RUN git clone --filter=tree:0 --depth=1 --no-checkout --sparse "${OSV_SCANNER_REPO}" . && \
  git sparse-checkout init --sparse-index --cone && \
  git sparse-checkout add cmd/osv-scanner internal pkg && \
  git checkout "${OSV_SCANNER_REF}"

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  go mod download && go mod verify

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  GOARCH=${TARGETARCH} GOOS=${TARGETOS} CGO_ENABLED=0 go build -a -ldflags="-s -w" -o osv-scanner ./cmd/osv-scanner/main.go

FROM --platform=${BUILDPLATFORM} golang:1.27.3@sha256:2920bb124f1a2e12a1d39f0cab756ae0068441ee1da7655dec7f5e374944095b AS builder

WORKDIR /src

# Use the toolchain specified in go.mod, or newer
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  go mod download && go mod verify

COPY cmd cmd
COPY internal internal

COPY --from=web-builder /src/internal/web/public /src/internal/web/public

ARG CUPDATE_VERSION="development build"
ARG TARGETARCH
ARG TARGETOS
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  GOARCH=${TARGETARCH} GOOS=${TARGETOS} CGO_ENABLED=0 go build -a -ldflags="-s -w -X 'main.Version=$CUPDATE_VERSION'" -o cupdate cmd/cupdate/*.go

FROM scratch AS export

COPY --from=builder /src/cupdate cupdate

FROM export

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=osv-scanner-builder /src/osv-scanner osv-scanner
COPY --from=web-builder /src/sbom.spdx.json /usr/local/share/sbom/frontend.spdx.json

ENV PATH=/

ENTRYPOINT ["cupdate"]
