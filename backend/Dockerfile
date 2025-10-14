# syntax=docker/dockerfile:1.4

FROM golang:1.22 AS base
WORKDIR /app

COPY go.mod ./
COPY go.sum ./go.sum
RUN --mount=type=cache,target=/go/pkg/mod \
    if [ -f go.sum ]; then go mod download; else echo "go.sum missing; run go mod tidy locally"; fi

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=base /app/bin/server /usr/local/bin/server
EXPOSE 8080
ENV PORT=8080
CMD ["/usr/local/bin/server"]
