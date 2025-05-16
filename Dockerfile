FROM --platform=${BUILDPLATFORM} node:24.0.2@sha256:8de41dd3ced665c49a1d7a0801f146fc88cd58ce53350fac59e5bd59c9ee8950 AS web-builder

WORKDIR /src

COPY .yarnrc.yml package.json yarn.lock .
COPY .yarn .yarn

RUN yarn install

COPY tsconfig.json vite.config.ts .
COPY web web

ARG CUPDATE_VERSION="development build"
RUN VITE_CUPDATE_VERSION="${CUPDATE_VERSION}" yarn build

FROM --platform=${BUILDPLATFORM} golang:1.24.2@sha256:30baaea08c5d1e858329c50f29fe381e9b7d7bced11a0f5f1f69a1504cdfbf5e AS builder

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
