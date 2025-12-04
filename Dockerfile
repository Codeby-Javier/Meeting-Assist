# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Runtime stage  
FROM python:3.9-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
RUN pip install --no-cache-dir \
    SpeechRecognition \
    PyMuPDF \
    Pillow \
    requests

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy necessary files
COPY templates ./templates
COPY static ./static
COPY scripts ./scripts
COPY .env .env

# Create upload directories
RUN mkdir -p uploads/audio uploads/documents uploads/exports

# Expose port
EXPOSE 8082

# Run the application
CMD ["./main"]
