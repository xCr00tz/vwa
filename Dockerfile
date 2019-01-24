## Docker setup for GoVWA


## setup GO & copying
FROM golang:1.8
RUN mkdir -p /go/src/vwa
WORKDIR /go/src/vwa
ADD . /go/src/vwa
RUN go get -v