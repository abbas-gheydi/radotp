
FROM golang:1.20-bullseye AS build

WORKDIR /

COPY . .

RUN ls -la

RUN go mod download && go build -o /radotp ./cmd/server/main.go 

FROM ubuntu:20.04 as radotp

LABEL maintainer="Abbas Gheydi <abbas.gheydi@gmail.com>"

RUN mkdir -p /etc/radotp

COPY radiusd.conf /etc/radotp

COPY --from=build /radotp /bin/radotp

RUN apt update && apt install -y openssl

ENTRYPOINT [ "/bin/radotp" ]

EXPOSE 8080 8081