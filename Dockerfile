FROM golang:1.12-alpine

RUN apk add git

COPY . /go/src/github.com/imdevinc/bl3_twitter_shift
WORKDIR /go/src/github.com/imdevinc/bl3_twitter_shift

ENV GO111MODULE=on

RUN go mod download && go mod verify

CMD go run cmd/watch/main.go