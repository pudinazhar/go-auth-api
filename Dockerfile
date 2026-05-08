# Stage 1: Build stage
FROM golang:1.26-alpine AS builder

# Install git (dibutuhkan untuk fetch library)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod dan sum files
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi (Target main.go yang sudah dipindah ke cmd/api)
RUN go build -o main ./cmd/main.go

# Stage 2: Run stage (Menggunakan image alpine yang sangat ringan)
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/main .
# Copy file .env
COPY .env .

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]