# Use the official Golang image as the base image
FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y ca-certificates openssl

ARG cert_location=/usr/local/share/ca-certificates

# Get certificate from "github.com"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Get certification from storage.googleapis.com
RUN openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/storage.googleapis.crt
# Update certificates
RUN update-ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

ENV CGO_ENABLED=0

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o gowebserver .

# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates file

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gowebserver .

# Copy the config.json file into the container
COPY config.json .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./gowebserver"]