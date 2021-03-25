package middleware

import (
	"context"
	"net/http"

	"github.com/corioders/gokit/constant"
	"github.com/corioders/gokit/web"
)

func Cors(origin string) web.Middleware {
	if !constant.IsProduction {
		origin = "*"
	}

	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			return handler(ctx, rw, r)
		}
	}
}
