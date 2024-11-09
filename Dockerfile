# Use an official Golang image as the base
FROM golang:1.21

# Install dependencies, including Chrome (or Chromium)
RUN apt-get update && apt-get install -y chromium

# Set the working directory inside the container
WORKDIR /app

# Copy your Go app code into the container
COPY . .

# Download Go module dependencies
RUN go mod download

# Build your Go app
RUN go build -o app .

# Expose the port your application listens on
EXPOSE 3000

# Run your app
CMD ["./app"]
