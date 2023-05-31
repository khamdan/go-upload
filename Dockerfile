# Use the official Golang base image
FROM golang:1.17-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o main .

# Use a minimal Alpine image as the final base image
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /app

# Copy the built executable from the build stage
COPY --from=build /app/main .

# Copy the frontend files
COPY ./fe /app/fe

RUN mkdir -p /app/uploads

# Expose the port your application listens on
EXPOSE 8000

# Run the Go application
CMD ["./main"]