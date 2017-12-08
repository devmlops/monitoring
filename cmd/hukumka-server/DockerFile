FROM golang:1.9.2-alpine

RUN mkdir -p "$GOPATH/src/github.com/wwwthomson/monitoring"
WORKDIR "$GOPATH/src/github.com/wwwthomson/monitoring"
COPY ./ ./

WORKDIR $GOPATH/bin
RUN go build github.com/wwwthomson/monitoring

WORKDIR "$GOPATH/src/github.com/wwwthomson/monitoring"
CMD ["/go/bin/hukumka-server"]
