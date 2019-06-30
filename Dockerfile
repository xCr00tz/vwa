## Docker setup for GoVWA


## setup GO & copying
FROM golang:1.8
RUN mkdir -p /go/src/vwa
WORKDIR /go/src/vwa
ADD . /go/src/vwa
RUN go get -v
## install additional package
## RUN go get github.com/jmoiron/sqlx
## RUN go get github.com/microcosm-cc/bluemonday
## RUN go get github.com/gorilla/csrf
## RUN go get github.com/justinas/nosurf
