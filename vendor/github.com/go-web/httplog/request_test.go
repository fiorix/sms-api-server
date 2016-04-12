package httplog

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUseXForwardedFor(t *testing.T) {
	want := "4.2.2.2:1905"
	var have string
	f := func(w http.ResponseWriter, r *http.Request) {
		have = r.RemoteAddr
	}
	r := &http.Request{
		Method:     "GET",
		URL:        &url.URL{Path: "/"},
		RemoteAddr: "[::1]:1905",
		Header: http.Header{
			"X-Forwarded-For": {"4.2.2.2"},
		},
	}
	w := &httptest.ResponseRecorder{}
	UseXForwardedFor(f).ServeHTTP(w, r)
	if have != want {
		t.Fatalf("Unexpected value for RemoteAddr. Want %q, have %q",
			want, have)
	}
}
