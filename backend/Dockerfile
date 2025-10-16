FROM golang:1.22 AS base
WORKDIR /app

COPY go.mod ./
COPY go.sum ./go.sum
RUN if [ -f go.sum ]; then go mod download; else echo "go.sum missing; run go mod tidy locally"; fi

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=base /app/bin/server /usr/local/bin/server
EXPOSE 8100
ENV PORT=8100
CMD ["/usr/local/bin/server"]
