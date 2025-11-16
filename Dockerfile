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
FROM --platform=${BUILDPLATFORM} golang:1.25.4@sha256:e68f6a00e88586577fafa4d9cefad1349c2be70d21244321321c407474ff9bf2 AS osv-scanner-builder

ARG TARGETARCH
ARG TARGETOS

ARG OSV_SCANNER_REPO="https://github.com/google/osv-scanner"
ARG OSV_SCANNER_REF="9e193da824d674ebaa946876a1905c7c3adbf114"

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

FROM --platform=${BUILDPLATFORM} golang:1.25.4@sha256:e68f6a00e88586577fafa4d9cefad1349c2be70d21244321321c407474ff9bf2 AS builder

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

ENV PATH=/

ENTRYPOINT ["cupdate"]
