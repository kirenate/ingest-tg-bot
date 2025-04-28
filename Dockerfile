FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ingestBot main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/ingest-bot /build/ingest-bot
CMD [". /ingest-bot"]
