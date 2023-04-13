# Building stage
FROM golang:1.17.3 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /go/src

COPY . /go/src

ENV TIMEZONE "Asia/Shanghai"

RUN go build -v -o _output/balancer ./

FROM reg.hrlyit.com/base/ubuntu:20.04 
COPY --from=builder /go/src/_output/balancer /opt/

ENTRYPOINT ["./opt/balancer"]
