# Start from the official Golang Alpine base image
FROM golang:1.21.6-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from your src directory
COPY src/go.mod src/go.sum ./

# Download any necessary dependencies
RUN go mod download

# Copy the rest of your source code from the src directory
COPY src/ .

# Build the Go app
# Ensure main.go or the relevant entry point is at the root of the src directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage: Use Alpine for the runtime image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
