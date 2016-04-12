package httpmux_test

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-web/httpmux"
)

func authHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok && u == "foobar" && p == "foobared" {
			ctx := httpmux.Context(r)
			ctx = context.WithValue(ctx, "user", u)
			httpmux.SetContext(ctx, r)
			next(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `realm="restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func Example() {
	// curl -i localhost:8080
	// curl -i -XPOST --basic -u foobar:foobared localhost:8080/auth/login
	root := httpmux.New()
	root.GET("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello, world\n")
	})
	auth := httpmux.New()
	{
		auth.UseFunc(authHandler)
		auth.POST("/login", func(w http.ResponseWriter, r *http.Request) {
			u := httpmux.Context(r).Value("user")
			fmt.Fprintln(w, "hello,", u)
		})
	}
	root.Append("/auth", auth)
	http.ListenAndServe(":8080", root)
}
