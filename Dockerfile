FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

RUN go build -o tunnel cmd/tunnel.go

# Create a new lightweight image
FROM alpine:3.18

WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=build /app/tunnel .

ENTRYPOINT ["./tunnel", "server"]
