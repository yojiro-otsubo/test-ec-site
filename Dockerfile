FROM golang:1.17.4-alpine

RUN apk update && \
    apk upgrade && \
    apk add git

WORKDIR /app

ADD ./src ./

RUN go get gopkg.in/ini.v1
RUN go get github.com/lib/pq
RUN go get -u github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/multitemplate
RUN go get github.com/gin-contrib/sessions
RUN go get github.com/utrack/gin-csrf
RUN go get -u github.com/stripe/stripe-go/v72
RUN go get github.com/stripe/stripe-go/v72/account
RUN go get github.com/stripe/stripe-go/v72/accountlink
RUN go get github.com/codenoid/gin-recaptcha