# Build Stage
FROM golang:1.21-alpine AS builder

# Set proxy to ensure downloads work
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

# Copy config files first
COPY go.mod go.sum ./

# Download dependencies (cache layer)
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage (Minimal Image)
FROM alpine:latest

# Install CA certificates for HTTPS APIs (AssemblyAI/OCRSpace)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy UI templates and assets
COPY templates/ ./templates/
COPY static/ ./static/

# Copy .env if exists (optional, simpler to use Railway Variables)
COPY .env* ./

# Create upload dirs
RUN mkdir -p uploads/audio uploads/documents uploads/exports

EXPOSE 8082

CMD ["./main"]
