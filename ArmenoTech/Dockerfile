# Use the official Golang image version 1.22.1 as the base image and name this stage "Builder"
FROM golang:1.22.2 AS Builder

# Set the working directory inside the container to /go/src/app
WORKDIR /go/src/app

# Copy all the files from the current directory (.) into the working directory inside the container
COPY . .

# Download module dependencies using Go modules
RUN go mod download

# Build the Go application with CGO disabled for a static binary and targeting Linux OS
RUN CGO_ENABLED=0 GOOS=linux go build -o telegrambot

# Start a new stage using the official Alpine Linux image as the base image
FROM alpine:latest

# Set the working directory inside the container to /app
WORKDIR /app

# Expose port 5432 for PostgreSQL
EXPOSE 5432

# Copy the built executable file "telegrambot" from the previous stage (Builder) into the /app directory of this stage
COPY --from=Builder /go/src/app/telegrambot .

# Command to run the application when the container starts
CMD ["./telegrambot"]
