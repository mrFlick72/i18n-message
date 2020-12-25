FROM golang:1.14.3

VOLUME /var/log/i18n-message
WORKDIR /app

COPY internal .
COPY configuration .
COPY api .
COPY main.go .
COPY go.mod .

RUN go get -d -v . && \
    go install -v . && \
    go build -o app

CMD ["./app"]