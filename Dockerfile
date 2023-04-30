# Set base image
FROM golang:1.19-alpine

# Set work dir
WORKDIR /app 

# Copy files and download Go modules
COPY go.* ./
RUN go mod download 

# Copy local code to container image
COPY . ./

# Build the binary app (discord bot)
RUN go build .

# Run the bot app
CMD ["./cordy"]