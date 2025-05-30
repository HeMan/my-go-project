# Use the official Go image for Go 1.24
FROM golang:1.24-alpine3.21 AS builder

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project, including go.mod and Go source files
COPY . ./

# Run tests to ensure the application is working correctly
RUN go test ./... -v

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o my-go-project

# Use a minimal base image for the final container
FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/my-go-project ./
COPY static ./static/

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./my-go-project"]
