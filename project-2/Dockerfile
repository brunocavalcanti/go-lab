FROM golang:1.10.3-alpine

WORKDIR /go/src/app
COPY . .
EXPOSE 3000
RUN apk update \
    && apk add sqlite sqlite-libs sqlite-dev  gcc musl-dev \
    && apk add socat \
    && apk add git

RUN go get -u github.com/mattn/go-sqlite3
RUN go get -d -v ./...
RUN go install -v ./...


CMD [ "app"]
