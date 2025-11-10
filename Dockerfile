FROM --platform=${BUILDPLATFORM} node:22.20.0@sha256:915acd9e9b885ead0c620e27e37c81b74c226e0e1c8177f37a60217b6eabb0d7 AS web-builder

WORKDIR /src

COPY .yarnrc.yml package.json yarn.lock .
COPY .yarn .yarn

RUN yarn install

COPY tsconfig.json vite.config.ts .
COPY web web

ARG CUPDATE_VERSION="development build"
RUN VITE_CUPDATE_VERSION="${CUPDATE_VERSION}" yarn build

# TODO: Download and install osv-scanner as an (optional) runtime dependency
# instead of including 100s of dependencies. Somehow include in SBOM...
FROM --platform=${BUILDPLATFORM} golang:1.25.3@sha256:6d4e5e74f47db00f7f24da5f53c1b4198ae46862a47395e30477365458347bf2 AS osv-scanner-builder

ARG TARGETARCH
ARG TARGETOS

ARG OSV_SCANNER_VERSION="v2.2.4"
ARG OSV_SCANNER_CHECKSUM_amd64="7702cd1e5d9f5059dd9570f4ad967f27d3c5f5391b371ec937b384c238177f55"
ARG OSV_SCANNER_CHECKSUM_arm64="94d1c520b30a7e28b0189b2a1dd24c7b08f41887186e8ae3f811067ec9ed7043"

SHELL ["/bin/bash", "-c"]

WORKDIR /src

RUN wget \
    -qO osv-scanner \
    "https://github.com/google/osv-scanner/releases/download/${OSV_SCANNER_VERSION}/osv-scanner_${TARGETOS}_${TARGETARCH}" && \
  OSV_SCANNER_CHECKSUM_VAR="OSV_SCANNER_CHECKSUM_${TARGETARCH}"; \
    echo "${!OSV_SCANNER_CHECKSUM_VAR}" "osv-scanner" | sha256sum --check --strict && \
  chmod +x osv-scanner

FROM --platform=${BUILDPLATFORM} golang:1.25.3@sha256:6d4e5e74f47db00f7f24da5f53c1b4198ae46862a47395e30477365458347bf2 AS builder

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
COPY --from=osv-scanner-builder /src/osv-scanner osv-scanner

ENV PATH=/

ENTRYPOINT ["cupdate"]
