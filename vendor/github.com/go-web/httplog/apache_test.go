package httplog

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-web/httpmux"
)

func TestApacheCommonFormat(t *testing.T) {
	var b bytes.Buffer
	l := log.New(&b, "", 0)
	mux := httpmux.New()
	mux.UseFunc(ApacheCommonFormat(l))
	mux.GET("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "foobar")
	}))
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/"},
	}
	w := &httptest.ResponseRecorder{}
	mux.ServeHTTP(w, r)
	// TODO: Test the format.
}

func TestApacheCombinedFormat(t *testing.T) {
	var b bytes.Buffer
	l := log.New(&b, "", 0)
	mux := httpmux.New()
	mux.UseFunc(ApacheCommonFormat(l))
	mux.GET("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "foobar")
	}))
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/"},
	}
	w := &httptest.ResponseRecorder{}
	mux.ServeHTTP(w, r)
	// TODO: Test the format.
}
