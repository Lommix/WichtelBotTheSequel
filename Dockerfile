FROM golang:1.23-alpine
RUN apk add --no-cache build-base

WORKDIR /usr/src/wichtelbot

# Pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
# Download and verify modules
RUN go mod download && go mod verify

# Copy contents to WORKDIR
COPY . .
# Build applications
RUN go build -v -o /usr/local/bin/ ./...

ENTRYPOINT ["wichtelbot"]
