FROM --platform=${BUILDPLATFORM} node:24.12.0@sha256:929c026d5a4e4a59685b3c1dbc1a8c3eb090aa95373d3a4fd668daa2493c8331 AS web-builder

WORKDIR /src

COPY .yarnrc.yml package.json yarn.lock .
COPY .yarn .yarn

RUN --mount=type=cache,target=node_modules  \
  yarn install --immutable

COPY tsconfig.json vite.config.ts .
COPY web web

ARG CUPDATE_VERSION="development build"
RUN --mount=type=cache,target=node_modules \
  VITE_CUPDATE_VERSION="${CUPDATE_VERSION}" yarn build

FROM --platform=${BUILDPLATFORM} golang:1.26.2@sha256:b54cbf583d390341599d7bcbc062425c081105cc5ef6d170ced98ef9d047c716 AS osv-scanner-builder

ARG TARGETARCH
ARG TARGETOS

ARG OSV_SCANNER_REPO="https://github.com/google/osv-scanner"
# v2.3.5 (2026-03-25)
ARG OSV_SCANNER_REF="30bcc134e23fbc35731021ee43ec433c483715d7"

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

FROM --platform=${BUILDPLATFORM} golang:1.26.2@sha256:b54cbf583d390341599d7bcbc062425c081105cc5ef6d170ced98ef9d047c716 AS builder

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
