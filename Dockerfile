# Builder Image
# -------------
FROM golang:1.14-alpine as builder
# System setup
RUN apk update && apk add git curl build-base autoconf automake libtool
# Install protoc
ENV PROTOBUF_URL https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protobuf-cpp-3.11.4.tar.gz
RUN curl -L -o /tmp/protobuf.tar.gz $PROTOBUF_URL
WORKDIR /tmp/
RUN tar xvzf protobuf.tar.gz
WORKDIR /tmp/protobuf-3.11.4
RUN mkdir /export
RUN ./autogen.sh && \
    ./configure --prefix=/export && \
    make -j 3 && \
    make check && \
    make install
# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go
RUN cp /go/bin/protoc-gen-go /export/bin/
# Export dependencies
RUN cp /usr/lib/libstdc++* /export/lib/
RUN cp /usr/lib/libgcc_s* /export/lib/


#   BASE IMAGE
# -------------
FROM golang:1.14-alpine as base
RUN apk update && apk add --no-cache git ca-certificates


#  DEV IMAGE
# -------------
FROM base as dev
COPY --from=builder /export /usr
# Add in air file monitor
RUN go get -u github.com/cosmtrek/air
WORKDIR /go/src/github.com/arkrozycki/go-pigeon
# Setup live reload
CMD air -c air.conf