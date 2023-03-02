############################
# STEP 1 build executable binary
############################
FROM golang:1.20.0-alpine3.17 AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v ./cmd/ 
# Build the binary.
RUN go build -o /go/bin/synse ./cmd/ 
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/synse /go/bin/synse
# Run the hello binary.
EXPOSE 8000
ENTRYPOINT ["/go/bin/synse"]