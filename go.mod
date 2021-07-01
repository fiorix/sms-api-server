module github.com/engagespark/sms-api-server

go 1.16

require (
	github.com/fiorix/go-smpp/v2 v2.0.0-20210507124410-a82a321cf1d8
	github.com/go-web/httplog v0.0.0-20160412232724-580d0d49f0d3
	github.com/go-web/httpmux v0.0.0-20160505070239-9e95425ee2c3 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
)

replace github.com/fiorix/go-smpp/v2 => github.com/alicebob/go-smpp/v2 v2.0.0
