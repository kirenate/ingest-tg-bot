FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ingest_bot main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/ingest_bot /build/ingest_bot
RUN chmod +x ingest_bot
CMD ["./ingest_bot"]
