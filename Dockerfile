FROM golang:1.8-alpine
MAINTAINER zhanglianx111@aliyun.com

RUN apk update && apk add gcc libc-dev curl
COPY . /go/src/github.com/zhanglianx111/clair-plus

RUN cd /go/src/github.com/zhanglianx111/clair-plus && go build -v -a -o /go/bin/clair-plus && cp /go/src/github.com/zhanglianx111/clair-plus/wait-for-postgres.sh / && rm -fr /go/src/github.com

EXPOSE 8080

CMD ["/go/bin/clair-plus"]

