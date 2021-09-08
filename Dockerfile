FROM golang as builder
ADD . /go/src/github.com/gugabfigueiredo/star-wars-api
WORKDIR /go/src/github.com/gugabfigueiredo/star-wars-api
RUN CGO_ENABLED=0 go build --mod=vendor -a -installsuffix cgo -o sw-api main.go

FROM golang:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/gugabfigueiredo/star-wars-api/sw-api .
CMD ./sw-api