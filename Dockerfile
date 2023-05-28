# Use a minimal base image
FROM golang:1.16-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the app
RUN go build -o tunnel cmd/tunnel.go

# Create a new lightweight image
FROM alpine:3.14

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=build /app/tunnel .

# Set the entrypoint for the container
ENTRYPOINT ["./tunnel"]
