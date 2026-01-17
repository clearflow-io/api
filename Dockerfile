# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 makes the binary statically linked (self-contained)
# GOOS=linux ensures it runs on Linux containers
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file if you want to include defaults (optional, environment variables in Render take precedence)
# COPY --from=builder /app/.env .

# Export the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
