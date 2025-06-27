FROM golang:1.23-alpine AS builder

# Sets the working directory inside the container to /app. All subsequent instructions will be run from this directory.
WORKDIR /app

# Copy go.mod and go.sum files to the working directory and download Go module dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the Go application for Linux.
# CGO_ENABLED=0 is important to make binary fully static.
# This command tells Go to build the main package located inside the ./cmd/api directory.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/api

# --- Final Stage ---
# Using a very small image because we only need the compiled result
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy ONLY compiled binary from 'builder' stage
COPY --from=builder /app/main .

# This command gives the operating system permission to run our program.
RUN chmod +x /app/main

# Expose port yang akan digunakan oleh aplikasi kita
EXPOSE 9000

# Command to run application when container starts
CMD ["/app/main"]