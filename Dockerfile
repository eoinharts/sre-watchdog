# ---- Build stage ----
    FROM golang:alpine AS builder

    WORKDIR /app
    
    # Copy dependency files first for better caching
    COPY go.mod go.sum ./
    
    RUN go mod download
    
    # Copy the application source
    COPY . .
    
    # Build Linux binary
    RUN CGO_ENABLED=0 GOOS=linux go build \
        -o /watchdog \
        ./cmd/watchdog
    
    
    # ---- Runtime stage ----
    FROM alpine:latest
    
    # Needed for HTTPS certificate verification
    RUN apk add --no-cache ca-certificates
    
    # Create non-root user
    RUN adduser -D -u 10001 watchdog
    
    # Copy compiled binary only
    COPY --from=builder /watchdog /usr/local/bin/watchdog
    
    USER watchdog
    
    EXPOSE 8080
    
    ENTRYPOINT ["/usr/local/bin/watchdog"]