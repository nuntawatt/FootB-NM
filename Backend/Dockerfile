# Multi-stage build for Go Backend
# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

# Install build dependencies
RUN apk add --no-cache ca-certificates git tzdata

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# Stage 2: Production Minimal Image
FROM alpine:3.20

WORKDIR /app

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy compiled binary from builder
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER appuser

ENV PORT=3000 \
    TZ=Asia/Bangkok

EXPOSE 3000

ENTRYPOINT ["/app/api"]
