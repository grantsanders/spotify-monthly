# Start from the official Golang base image
FROM golang:1.20.13-bullseye as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files (if your project is using Go modules)
# This is to ensure that dependencies are cached unless these files change
COPY go.mod go.sum ./

# Download any necessary dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch for a smaller, final image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
