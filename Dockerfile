FROM golang:1.19-alpine AS build

RUN mkdir /app

ADD . /app
WORKDIR /app

RUN go mod download

RUN go build -o /notionbot cmd/main.go

FROM alpine:latest AS deploy

WORKDIR /

COPY --from=build /notionbot /notionbot


ENTRYPOINT ["/notionbot"]