# #########################################
#   BASE IMAGE
FROM golang:1.14-alpine as base
RUN apk update && apk add --no-cache git

# #########################################
#  DEV IMAGE
FROM base as dev
# ADD in air file monitor
RUN go get -u github.com/cosmtrek/air
# COPY ./air.conf /tmp/air.conf
WORKDIR /go/src/github.com/arkrozycki/go-pigeon
# SETUP LIVE RELOAD
# turns up both inbound and outbound; configured in air.conf
CMD air -c air.conf