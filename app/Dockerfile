FROM golang:1-alpine as builder

ENV GO111MODULE=on

RUN apk add --no-cache git upx musl-dev curl
RUN go get github.com/pwaller/goupx

ADD app/ /root/app/

WORKDIR /root/app/
RUN go build -ldflags "-s -w" server.go
RUN goupx server

FROM alpine:latest as runner
WORKDIR /root/app/
RUN apk add --no-cache libc6-compat ca-certificates
COPY --from=builder /root/app/server /root/app/
CMD ["/bin/sh", "-c", "sleep 10 && exec /root/app/server"]
