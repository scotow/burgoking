##################################
# STEP 1 build executable binary #
##################################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
COPY . $GOPATH/src/github.com/scotow/burgoking

# Move to command directory.
WORKDIR $GOPATH/src/github.com/scotow/burgoking/cmd/web_original

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/burgoking

##############################
# STEP 2 build a small image #
##############################
FROM scratch

# Copy our static executable and static files.
COPY --from=builder /go/bin/burgoking /burgoking
COPY cmd/web_original/static /static

# Copy SSL certificates for HTTPS connections.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy locale data.
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip

# Run the hello binary.
ENTRYPOINT ["/burgoking"]