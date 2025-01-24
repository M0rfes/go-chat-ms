FROM golang:1.23-alpine as builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code
COPY . .

# Build the application
RUN go build -o main

# Deploy stage
FROM alpine:latest

WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/main .

# Copy templates
COPY templates ./templates

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
