FROM golang:1.16 as builder
WORKDIR /build
COPY ./ ./
RUN make build

FROM debian:stable-slim
#RUN apt update && apt install -y ca-certificates
COPY --from=0 /build/sms-api-server  /
COPY --from=0 /build/index.html  /pub/index.html

EXPOSE 8080
ENTRYPOINT ["/sms-api-server", "-public", "/pub"]
