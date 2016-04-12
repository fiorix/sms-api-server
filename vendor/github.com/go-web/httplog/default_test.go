package httplog

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-web/httpmux"
)

func TestDefaultFormat(t *testing.T) {
	var b bytes.Buffer
	l := log.New(&b, "", 0)
	mux := httpmux.New()
	mux.UseFunc(DefaultFormat(l))
	mux.GET("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "foobar")
	}))
	r := &http.Request{
		Proto:  "HTTP/1.1",
		Method: "GET",
		URL:    &url.URL{Path: "/"},
	}
	w := &httptest.ResponseRecorder{}
	mux.ServeHTTP(w, r)
	// TODO: Test the format.
	if !strings.HasPrefix(b.String(), "HTTP/1.1 404") {
		t.Fatalf("Unexpected data. Want HTTP/1.1 404, have %q", b.String())
	}
}
