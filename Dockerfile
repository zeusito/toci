# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o myapp ./cmd/main.go

RUN go test -v ./...

# Stage 3: Start fresh from a smaller image
FROM gcr.io/distroless/static

# Copy the binary from the builder stage
COPY --from=builder /app/myapp /app/myapp

# Set working directory
WORKDIR /app

# Expose the port your application listens on
EXPOSE 8080

# Set the entrypoint for the container
ENTRYPOINT ["./myapp"]
