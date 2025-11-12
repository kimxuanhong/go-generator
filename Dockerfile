# Multi-stage Dockerfile for building and running the Go application
# Builder stage
FROM golang:1.20 AS builder

WORKDIR /src

# Copy go.mod/go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources
COPY . .

# Build a statically linked binary
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -trimpath -ldflags "-s -w" -o /app/goex ./main.go

# Final stage: small runtime image
FROM alpine:3.18 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/goex ./goex

# Non-root user for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Application typically listens on 8080; update if your app uses a different port
EXPOSE 8080

ENTRYPOINT ["/app/goex"]
