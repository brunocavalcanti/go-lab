FROM golang:1.8-alpine

WORKDIR /go/src/app
COPY . .
EXPOSE 3000
RUN apk add --no-cache git -y
RUN go get -d -v ./...
RUN go install -v ./...


CMD [ "app"]