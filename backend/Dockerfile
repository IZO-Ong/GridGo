# Stage 1: Build binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first to cache layers
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# lightweight image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]