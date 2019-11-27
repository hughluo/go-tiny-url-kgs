FROM golang:1.13.4

WORKDIR /go/src/app
COPY . .

RUN go get -u ./...
RUN go build

CMD ["app"]


