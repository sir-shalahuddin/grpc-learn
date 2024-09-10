# Use the official Go image as the base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Move to the directory containing the main.go file
WORKDIR /app/cmd

# Build the Go application
RUN go build -o main .

# Expose the application port
EXPOSE 3000
EXPOSE 3001


# Run the application
CMD ["./main"]
