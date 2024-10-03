# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Run unit tests
RUN go test -v ./...

# Move to the directory containing the main.go file
WORKDIR /app/cmd

# Build the Go application
RUN go build -o main .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install necessary dependencies for running the application
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/cmd/main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
