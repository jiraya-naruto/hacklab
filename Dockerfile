# Start from a lightweight official Golang image
FROM golang:1.21-alpine AS builder

# Install dependencies for building Go application
RUN apk add --no-cache --update \
    gcc \
    g++ \
    libtool \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o app .

# Start a new, smaller image for deployment
FROM alpine:latest

# Install Chrome runtime dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    libx11 \
    libxcomposite \
    libxdamage \
    libxi \
    mesa-gl

# Copy the compiled application from the builder
COPY --from=builder /app/app /app

WORKDIR /app

# Expose the application’s port (change this if your app uses a different port)
EXPOSE 8080

# Run the application
CMD ["./app"]