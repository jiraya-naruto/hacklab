# Use an official Golang image as the base
FROM golang:1.21

# Install dependencies, including Chrome (or Chromium)
RUN apt-get update && apt-get install -y chromium

# Set the working directory inside the container
WORKDIR /app

# Copy your Go module files and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of your Go app code into the container
COPY . .

# Ensure main.go is in the /app directory
RUN ls -la /app  # This will print the directory contents to help with debugging

# Build your Go app
RUN go build -o app main.go  # Specify main.go to make sure it targets the correct file

# Expose the port your application listens on
EXPOSE 3000

# Run your app
CMD ["./app"]
