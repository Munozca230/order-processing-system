# Multi-stage build for optimized production image
FROM golang:1.22-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o product-api-v2 \
    ./cmd/server

# Production stage
FROM alpine:3.19

# Install ca-certificates for HTTPS requests and timezone data
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user for security
RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/product-api-v2 .

# Change ownership to non-root user
RUN chown appuser:appgroup /app/product-api-v2 && \
    chmod +x /app/product-api-v2

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV PORT=8080 \
    ENVIRONMENT=production \
    LOG_LEVEL=info \
    LOG_FORMAT=json

# Run the application
ENTRYPOINT ["./product-api-v2"]