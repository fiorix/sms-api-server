package httplog

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/context"

	"github.com/go-web/httpmux"
)

func TestError(t *testing.T) {
	var ctx context.Context
	mux := httpmux.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		Error(r, "hello", "world")
		ctx = httpmux.Context(r)
	})
	r := &http.Request{Method: "GET", URL: &url.URL{Path: "/"}}
	mux.ServeHTTP(&httptest.ResponseRecorder{}, r)
	if ctx.Value(ErrorID) != "helloworld" {
		t.Fatalf("Unexpected value. Want \"helloworld\", have %v", ctx.Value("error"))
	}
}

func TestErrorf(t *testing.T) {
	var ctx context.Context
	mux := httpmux.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		Errorf(r, "hello, world")
		ctx = httpmux.Context(r)
	})
	r := &http.Request{Method: "GET", URL: &url.URL{Path: "/"}}
	mux.ServeHTTP(&httptest.ResponseRecorder{}, r)
	if ctx.Value(ErrorID) != "hello, world" {
		t.Fatalf("Unexpected value. Want \"hello, world\", have %v", ctx.Value("error"))
	}
}

func TestErrorln(t *testing.T) {
	var ctx context.Context
	mux := httpmux.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		Errorln(r, "hello", "world")
		ctx = httpmux.Context(r)
	})
	r := &http.Request{Method: "GET", URL: &url.URL{Path: "/"}}
	mux.ServeHTTP(&httptest.ResponseRecorder{}, r)
	if ctx.Value(ErrorID) != "hello world\n" {
		t.Fatalf("Unexpected value. Want \"hello world\\n\", have %v", ctx.Value("error"))
	}
}
