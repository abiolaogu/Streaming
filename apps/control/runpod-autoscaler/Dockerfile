FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum if they exist
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source code
COPY *.go ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o runpod-autoscaler -ldflags='-w -s' .

FROM alpine:latest

RUN apk --no-cache add ca-certificates openssh-client

WORKDIR /app

COPY --from=builder /build/runpod-autoscaler .

ENTRYPOINT ["/app/runpod-autoscaler"]

