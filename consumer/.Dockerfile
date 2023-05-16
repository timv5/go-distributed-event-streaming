# Use the official Go image for version 1.18 as the base image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app/consumer

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN go build -o app

# Expose a port (change it to the port your Go app listens on)
EXPOSE 8080

# Specify the command to run the binary within the container
CMD ["./app"]
