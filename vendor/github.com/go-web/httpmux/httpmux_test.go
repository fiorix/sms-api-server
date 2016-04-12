package httpmux

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/context"
)

func TestHandler(t *testing.T) {
	mux := New()
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.DELETE("/", f)
	mux.GET("/", f)
	mux.HEAD("/", f)
	mux.OPTIONS("/", f)
	mux.PATCH("/:arg", f)
	mux.POST("/:arg", f)
	mux.PUT("/:arg", f)
	for i, method := range []string{
		"DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT",
	} {
		r := &http.Request{
			Method: method,
			URL:    &url.URL{Path: "/"},
		}
		switch method {
		case "PATCH", "POST", "PUT":
			r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{1}))
			r.URL.Path += "arg"
		}
		w := &httptest.ResponseRecorder{}
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("Test %d: unexpected status code. Want %d, have %d",
				i, http.StatusOK, w.Code)
		}
	}
}

func TestSubtree(t *testing.T) {
	root := New()
	root.UseFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Hello", "world")
			next(w, r)
		}
	})
	c := DefaultConfig
	c.Prefix = "/ignore-me"
	c.UseFunc(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Hello", h.Get("X-Hello")+"z")
			next(w, r)
		}
	})
	tree := NewHandler(&c)
	tree.GET("/foobar", func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("X-Hello") == "worldz" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	})
	root.Append("/test", tree)
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test/foobar"},
	}
	w := &httptest.ResponseRecorder{}
	root.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Unexpected status code. Want %d, have %d",
			http.StatusOK, w.Code)
	}

}

func testmw(want int) Middleware {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			p := Params(r).ByName("opt")
			if p != "foobar" {
				http.Error(w, "missing parameter: foobar",
					http.StatusNotFound)
				return
			}
			ctx := Context(r)
			have, _ := ctx.Value("v").(int)
			if want != have {
				m := fmt.Sprintf("want=%d have=%d", want, have)
				http.Error(w, m, http.StatusNotFound)
				return
			}
			have++
			ctx = context.WithValue(ctx, "v", have)
			SetContext(ctx, r)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}

func TestMiddleware(t *testing.T) {
	root := New()
	root.Use(testmw(0))
	root.Use(testmw(1))
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	root.GET("/:opt", http.HandlerFunc(f))
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/foobar"},
	}
	w := &httptest.ResponseRecorder{Body: &bytes.Buffer{}}
	root.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Middleware chain is broken: %s", w.Body.Bytes())
	}
}

func TestMiddlewareSubtree(t *testing.T) {
	root := New()
	root.Use(testmw(0))
	root.Use(testmw(1))
	subtree := New()
	subtree.Use(testmw(2))
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	subtree.GET("/:opt", http.HandlerFunc(f))
	root.Append("/a", subtree)
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/a/foobar"},
	}
	w := &httptest.ResponseRecorder{Body: &bytes.Buffer{}}
	root.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Middleware chain is broken: %s", w.Body.Bytes())
	}
}

func TestServeFiles(t *testing.T) {
	i := 0
	p := map[string]struct{ Dir, URL string }{
		"/*filepath":        {".", "/httpmux.go"},
		"/foobar/*filepath": {".", "/foobar/httpmux.go"},
	}
	for pattern, cfg := range p {
		mux := New()
		mux.ServeFiles(pattern, http.Dir(cfg.Dir))
		w := &httptest.ResponseRecorder{}
		r := &http.Request{
			Method: "GET",
			URL:    &url.URL{Path: cfg.URL},
		}
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("Test %d: Unexpected status. Want %d, have %d",
				i, http.StatusOK, w.Code)
		}
		i++
	}
}
