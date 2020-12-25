FROM golang:1.14.3

VOLUME /var/log/i18n-message
WORKDIR /go/src/app

COPY docker .

RUN go get -d -v ./...  && \
    go install -v ./... && \
    go build -o main

CMD ["./cmd/app/main"]