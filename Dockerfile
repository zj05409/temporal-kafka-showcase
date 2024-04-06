# Start from the Go version specified in go.mod
FROM golang:1.22.1

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o starter_main ./starter/main.go
RUN go build -o worker_main ./worker/main.go
RUN chmod +x ./main.sh

ENTRYPOINT ["sh", "-c", "./main.sh"]