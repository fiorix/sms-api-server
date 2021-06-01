This is a fork of https://github.com/fiorix/sms-api-server


# HTTP API for sending SMS via SMPP

The sms-api-server is a web server that connects to an SMSC via
SMPP v3.4. It provides HTTP and WebSocket APIs for sending short
messages and querying their status (when supported by the SMSC).
It supports sending delivery receipts via WebSockets or
[Server-Sent Events](http://www.w3schools.com/html/html5_serversentevents.asp).

[![Build Status](https://secure.travis-ci.org/fiorix/sms-api-server.png)](https://travis-ci.org/fiorix/sms-api-server)

## Running

Developers can build the source code and run. For everyone
else, the easiest way to run is via Docker:

	docker run --rm -i -t fiorix/sms-api-server [--help]

See [this link](https://hub.docker.com/r/fiorix/sms-api-server/)
on the Docker hub for details.

## Usage

With the server running, send a message:

	curl localhost:8080/v1/send -X POST -F src=bart -F dst=lisa -F text=hi

In case of success, the server returns a JSON document containing a
message ID that can be used for querying its delivery status later.
This functionality is not always available on some SMSCs.

Example:

	curl "localhost:8080/v1/query?src=bart&message_id=1234"

To watch for incoming SMS, or delivery receipts:

	curl localhost:8080/v1/sse

This is the Server-Sent Events (SSE) endpoint that deliver messages
as events, as they arrive on the server.

## Send parameters

The `/v1/send` endpoint supports the following parameters:

- src: number of sender (optional)
- dst: number of recipient
- text: text message, encoded as UTF-8
- enc: text encoding for short message: `latin1` or `ucs2` (optional)
- register: register for delivery: `final` or `failure` (optional)

If an encoding is not provided, data is sent as a binary blob and may
not display correctly on some devices.

For special characters, try:

	curl localhost:8080/v1/send -X POST -F dst=foobar -F enc=ucs2 -F text="é nóis"

## WebSocket API

This server provides two websocket APIs:

- One for sending messages and querying for message status
- One for sending delivery receipts

See [index.html](./index.html) for details.

Have fun!
