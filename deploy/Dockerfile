
FROM golang:1.17 AS build

WORKDIR /

COPY . .

RUN ls -la

RUN go mod download && go build -o /radotp .

FROM ubuntu:20.04 as radotp

LABEL maintainer="Abbas Gheydi <abbas.gheydi@gmail.com>"

RUN mkdir -p /etc/radotp

COPY radiusd.conf /etc/radotp

COPY --from=build /radotp /bin/radotp

ENTRYPOINT [ "/bin/radotp" ]

EXPOSE 8080 8081