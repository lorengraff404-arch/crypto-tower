# Build Stage
FROM golang:1.23-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git build-base

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# -ldflags="-w -s" reduces binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api

# Final Stage
FROM alpine:3.19

WORKDIR /app

# Install certificates for external HTTPS calls (Blockchain/Auth)
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/api .

# Copy configuration
COPY --from=builder /app/config ./config

# Copy static assets for Admin Dashboard
COPY --from=builder /app/admin-dashboard ./admin-dashboard

# Expose port
EXPOSE 8080

# Environment variables (Defaults)
ENV PORT=8080
ENV APP_ENV=production

# Run the binary
CMD ["./api"]
