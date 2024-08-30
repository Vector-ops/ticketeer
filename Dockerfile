FROM golang:1.23.0-alpine3.19

# Set the working directory
WORKDIR /src/app

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .


# Command for running the application using air
CMD ["air", "-c", ".air.toml"]
