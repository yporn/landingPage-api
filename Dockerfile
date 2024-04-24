# Use Go image as a builder
FROM --platform=linux/amd64 golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install GCC
RUN apk add --no-cache gcc libc-dev

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install migrate tool
RUN apk --no-cache add curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate \
    && apk del curl

# Clean up unused dependencies
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use lightweight Alpine image for the final build stage
FROM alpine:3.16

# Set working directory
WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

# Copy migrate binary and migrations directory from builder stage
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/pkg/databases/migrations /app/pkg/databases/migrations

# Expose port 3000
EXPOSE 3000

# Define command to run the application
CMD ["/app/main"]
