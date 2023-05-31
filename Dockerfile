FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

RUN go build -o auger cmd/auger.go

# Create a new lightweight image
FROM alpine:3.18

WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=build /app/auger .

ENTRYPOINT ["./auger", "server"]
