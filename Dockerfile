FROM golang:1.6

ADD . /go/src/github.com/fiorix/sms-api-server
COPY index.html /pub/index.html
RUN GO15VENDOREXPERIMENT=1 go install github.com/fiorix/sms-api-server

EXPOSE 8080
ENTRYPOINT ["/go/bin/sms-api-server", "-public", "/pub"]
