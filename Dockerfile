FROM node:22 AS web-builder

WORKDIR /src

COPY .yarnrc.yml package.json yarn.lock .
COPY .yarn .yarn

RUN yarn install

COPY eslint.config.mjs postcss.config.mjs tailwind.config.js tsconfig.json vite.config.ts .
COPY web web

RUN yarn build

FROM golang:1.23 AS builder

WORKDIR /src

COPY go.mod go.sum .

RUN go mod download

COPY cmd cmd
COPY internal internal

COPY --from=web-builder /src/internal/web/public /src/internal/web/public

RUN CGO_ENABLED=0 go build -a -ldflags="-s -w" -o cupdate cmd/cupdate/*.go

FROM scratch AS export

COPY --from=builder /src/cupdate cupdate

FROM export

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV PATH=/

ENTRYPOINT ["cupdate"]