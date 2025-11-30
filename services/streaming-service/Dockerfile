FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum
COPY go.mod go.sum ./

# Copy shared packages
COPY packages/common-go ./packages/common-go
COPY packages/proto ./packages/proto

# Copy service source
COPY services/streaming-service ./services/streaming-service

# Build the service
WORKDIR /app/services/streaming-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o streaming-service ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/streaming-service/streaming-service .

EXPOSE 8080

CMD ["./streaming-service"]
