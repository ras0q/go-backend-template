# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod

RUN \
  --mount=type=cache,target=${GOCACHE} \
  --mount=type=cache,target=${GOMODCACHE} \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  go mod download

RUN \
  --mount=type=cache,target=${GOCACHE} \
  --mount=type=cache,target=${GOMODCACHE} \
  --mount=type=bind,target=. \
  go build -o /usr/bin/server ./cmd/server/main.go

# use `debug-nonroot` for debug shell access
FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /app

COPY --from=builder /usr/bin/server /usr/bin/server

CMD ["/usr/bin/server"]
