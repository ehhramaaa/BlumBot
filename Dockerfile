################
# BUILD BINARY #
################

FROM golang:latest AS builder

WORKDIR /app
COPY . .

# Download dependencies and verify
RUN go mod download
RUN go mod verify

# Build the Go binary with the specified flags
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/BlumBot .

#####################
# FINAL IMAGE #
#####################

FROM alpine:3.16

# Set working directory in the final image
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/BlumBot /app/

# Copy the configs directory from the build context
COPY . /app/

# Ensure the binary has execution permissions
RUN chmod +x /app/BlumBot

# Set the entrypoint for the container
CMD ["./BlumBot"]