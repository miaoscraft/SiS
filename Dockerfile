FROM golang:1.24 AS builder

COPY . /src

WORKDIR /src

RUN go build -o sis

FROM ubuntu:latest
LABEL authors="BaiMeow"

COPY --from=builder /src/sis /app/sis

WORKDIR /app

ENTRYPOINT ["sis"]