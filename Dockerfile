FROM golang:1.8

WORKDIR /

COPY ./ /go/src/github.com/blackducksoftware/canary/

WORKDIR /go/src/github.com/blackducksoftware/canary/

RUN go build cmd/main.go

CMD ./main
