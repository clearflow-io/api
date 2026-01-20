# Stage 1: Build
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Final Image
FROM alpine:latest  

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

ENV PORT=8080
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./main"]
