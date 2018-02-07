FROM golang:latest

RUN mkdir -p /go/src/github.com/iZIVer/imagemaker

WORKDIR /go/src/github.com/iZIVer/imagemaker

COPY . /go/src/github.com/iZIVer/imagemaker

RUN go-wrapper download

RUN go-wrapper install

CMD ["go-wrapper", "run", "-web"]

EXPOSE 8000