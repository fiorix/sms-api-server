FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/sms-api-server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/sms-api-server /usr/local/bin/sms-api-server
COPY index.html /pub/index.html
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/sms-api-server", "-public", "/pub"]
