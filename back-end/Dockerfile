# https://docs.docker.com/build/building/multi-stage/
# I just didn't add the alpine linux because that would be more overhead that wouldn't payoff
FROM golang:1.23.4

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
