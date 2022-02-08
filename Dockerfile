# syntax=docker/dockerfile:1

FROM golang:1.17.6-alpine
WORKDIR /build


RUN apk add build-base
# Install dependencies
# Thanks to @montanaflynn
# https://github.com/montanaflynn/golang-docker-cache
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

COPY . .

RUN go build -o /giftDetester
ENTRYPOINT [ "/giftDetester" ]