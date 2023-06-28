# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Gnod in a stock Go builder container
FROM golang:1.18-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /go-neworder/
COPY go.sum /go-neworder/
RUN cd /go-neworder && go mod download

ADD . /go-neworder
RUN cd /go-neworder && go run build/ci.go install -static ./cmd/gnod

# Pull Gnod into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-neworder/build/bin/gnod /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["gnod"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
