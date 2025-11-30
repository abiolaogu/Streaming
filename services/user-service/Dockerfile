FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy root go.mod and go.sum
COPY go.mod go.sum ./

# Copy shared packages
COPY packages/common-go ./packages/common-go
COPY packages/proto ./packages/proto

# Copy service source
COPY services/user-service ./services/user-service

# Build the service
WORKDIR /app/services/user-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/user-service/user-service .

EXPOSE 8080

CMD ["./user-service"]
