FROM golang:1.23-alpine as builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application code
COPY . .

# Build the binary
RUN go build -o main .

# Deploy stage
FROM alpine:latest

WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
