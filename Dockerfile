# syntax=docker/dockerfile:1

FROM node:22 AS frontend-builder

WORKDIR /app

RUN \
  --mount=type=bind,source=frontend/app-ui/package.json,target=package.json \
  --mount=type=bind,source=frontend/app-ui/package-lock.json,target=package-lock.json \
  npm ci

COPY ./frontend/app-ui ./
RUN npm run build
RUN ls -al

FROM golang:1.25 AS builder

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

COPY --from=frontend-builder /app/dist /tmp/dist
RUN \
  --mount=type=cache,target=${GOCACHE} \
  --mount=type=cache,target=${GOMODCACHE} \
  --mount=type=bind,target=.,readwrite \
  cp -r /tmp/dist /app/frontend/app-ui/dist \
  && go build -o /usr/bin/server ./main.go

# use `debug-nonroot` for debug shell access
FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /app

COPY --from=builder /usr/bin/server /usr/bin/server

CMD ["/usr/bin/server"]
