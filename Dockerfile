# Use the official Go image for Go 1.24
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project, including go.mod and Go source files
COPY . ./

# Generate go.sum and resolve dependencies
RUN go mod tidy

# Run tests to ensure the application is working correctly
RUN go test ./... -v

# Build the Go application
RUN go build -o my-go-project

# Use a minimal base image for the final container
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Install PostgreSQL client tools in the final container
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder stage
COPY --from=builder /app/my-go-project ./

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./my-go-project"]
