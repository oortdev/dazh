# ----------------------------
# Stage 1: Build the Go binary
# ----------------------------
FROM golang:1.21-alpine AS builder

# Install git and bash for module support
RUN apk add --no-cache git bash

# Set working directory
WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dashboard main.go

# ----------------------------
# Stage 2: Create minimal image
# ----------------------------
FROM gcr.io/distroless/base

# Set working directory inside container
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/dashboard .

# Copy templates, static files, and config
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 3000

# Run the application
CMD ["./dashboard"]

