FROM golang:1.18.4

VOLUME /var/log/onlyone-portal/logs
WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...  && \
    go install -v ./... && \
    go build -o app

CMD ["./app"]