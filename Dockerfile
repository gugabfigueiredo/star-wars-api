FROM golang as builder
ADD . /go/src/github.com/gugabfigueiredo/star-wars-api
WORKDIR /go/src/github.com/gugabfigueiredo/star-wars-api
RUN CGO_ENABLED=0 go build --mod=vendor -a -installsuffix cgo -o sw-api main.go