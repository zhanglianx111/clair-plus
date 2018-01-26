FROM golang:1.8-alpine
MAINTAINER zhanglianx111@aliyun.com

RUN apk update && apk add gcc libc-dev curl
COPY . /go/src/github.com/zhanglianx111/clair-plus

RUN cd /go/src/github.com/zhanglianx111/clair-plus && \
    go build -v -a -o /bin/clair-plus && \
    rm -fr /go/src/github.com

EXPOSE 8080

CMD ["clair-plus"]

