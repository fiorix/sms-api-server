package httplog

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseRecorder(t *testing.T) {
	w := NewResponseWriter(&httptest.ResponseRecorder{})
	w.Header().Set("X-test", "foobar")
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "hello, world")
	switch {
	case w.Code() != http.StatusNotFound:
		t.Fatalf("Unexpected status code. Want %d, have %d",
			http.StatusNotFound, w.Code())
	case w.Bytes() != 12:
		t.Fatalf("Unexpected # of bytes. Want 12, have %d", w.Bytes())
	}
}
