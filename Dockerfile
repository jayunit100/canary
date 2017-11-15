FROM golang:1.8

WORKDIR /
# /Users/jayunit100/work/go/src/bitbucket.org/bdsengineering/hub-sidecar
COPY ./ /go/src/bitbucket.org/bdsengineering/hub-sidecar

WORKDIR /go/src/bitbucket.org/bdsengineering/hub-sidecar
RUN ls -altrh
RUN go build cmd/main.go

CMD ./main
