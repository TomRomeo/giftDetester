# syntax=docker/dockerfile:1

FROM golang:1.17.6-alpine
WORKDIR /build

COPY . .

RUN go get .

RUN go build -o /giftDetester
ENTRYPOINT [ "/giftDetester" ]