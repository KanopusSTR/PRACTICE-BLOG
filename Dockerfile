# use official Golang image
FROM golang:1.22

# set working directory
WORKDIR /server

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o server ./cmd/app/

# Run the executable
CMD ["./server"]
