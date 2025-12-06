# Use official Go image
FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build
RUN go build -o main .

# Expose port
EXPOSE 8082

# Run
CMD ["./main"]
