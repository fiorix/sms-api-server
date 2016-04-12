package httplog_test

import (
	"log"
	"net/http"
	"os"

	"github.com/go-web/httplog"
	"github.com/go-web/httpmux"
)

func ExampleLog() {
	logger := log.New(os.Stdout, "[go-web] ", log.LstdFlags)
	mux := httpmux.New()
	mux.UseFunc(httplog.DefaultFormat(logger))
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		httplog.Errorf(r, "Todos são manos, eeei. Todos são hu-manos.")
	})
	http.ListenAndServe(":8080", mux)
}
