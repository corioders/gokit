package middleware

import (
	"context"
	"net/http"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/log"
	"github.com/corioders/gokit/web"
)

// Errors middelware catches errors and recovers from panics.
func Errors(logger log.Logger) web.Middleware {
	logger = logger.Child("Errors middleware")
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
			defer func() {
				r := recover()

				if r != nil {
					if err, ok := r.(error); ok {
						logger.Error(errors.WithMessage(err, "Panic"))
						return
					}

					logger.Error("Panic:", r)
				}
			}()

			err := handler(ctx, rw, r)

			if err != nil {
				logger.Error(errors.WithMessage(err, "Error"))
			}

			return nil
		}
	}
}
