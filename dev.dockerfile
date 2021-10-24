FROM golang:1.16-alpine

ENV GO111MODULE=on

RUN go get github.com/cortesi/modd/cmd/modd

WORKDIR /app

EXPOSE 8080

ENTRYPOINT /bin/sh -c "modd -f ./tools/modd.conf"
