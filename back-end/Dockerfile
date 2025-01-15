# https://docs.docker.com/build/building/multi-stage/ Used as reference, although I don't have another container (tester)

# Stage 1
# This one will build and run the application
FROM golang:1.23.4 AS builder

# Set working directory
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the app and send it to /bin/
RUN go build -o /bin/sipub-test ./cmd/sipub-test/main.go

# Expose the port
EXPOSE 8080

CMD ["/bin/sipub-test"]
