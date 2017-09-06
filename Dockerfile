FROM golang:1.9-alpine

WORKDIR /www

RUN apk add --no-cache --update git
RUN go get github.com/phinexdaz/ipapk-server
RUN go build

CMD ["./ipapk-server"]