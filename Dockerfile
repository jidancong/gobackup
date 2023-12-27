FROM golang:1.21.5 AS builder

WORKDIR /app

COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go build -o gobackup

FROM ubuntu:22.04

WORKDIR /app

RUN apt update \
    && apt install mysql-client -y \
    && apt install postgresql-client -y

COPY --from=builder /app/gobackup .

ENTRYPOINT ["./gobackup"]