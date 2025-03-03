FROM --platform=${BUILDPLATFORM} node:22.14.0@sha256:cfef4432ab2901fd6ab2cb05b177d3c6f8a7f48cb22ad9d7ae28bb6aa5f8b471 AS web-builder

WORKDIR /src

COPY .yarnrc.yml package.json yarn.lock .
COPY .yarn .yarn

RUN yarn install

COPY tailwind.config.js tsconfig.json vite.config.ts .
COPY web web

ARG CUPDATE_VERSION="development build"
RUN VITE_CUPDATE_VERSION="${CUPDATE_VERSION}" yarn build

FROM --platform=${BUILDPLATFORM} golang:1.24.0@sha256:3f7444391c51a11a039bf0359ee81cc64e663c17d787ad0e637a4de1a3f62a71 AS builder

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
