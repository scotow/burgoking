################################
# STEP 1 build executable binary
################################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/scotow/burgoking/cmd/web
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/burgoking

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable and static files.
COPY --from=builder /go/bin/burgoking /burgoking
COPY static /static

# Copy SSL certificates for HTTPS connections.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


# Run the hello binary.
ENTRYPOINT ["/burgoking"]