package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/corioders/gokit/log"

	"github.com/dimfeld/httptreemux"
)

type Handler func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error

type WebRouter struct {
	logger log.Logger

	router     *httptreemux.ContextMux
	middleware []Middleware
}

func NewRouter(logger log.Logger, middleware ...Middleware) *WebRouter {
	return &WebRouter{
		logger: logger,

		router:     httptreemux.NewContextMux(),
		middleware: middleware,
	}
}

// Handle registers handler on specified path and method
func (w *WebRouter) Handle(method string, path string, handler Handler, middleware ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(middleware, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(w.middleware, handler)

	h := func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := handler(ctx, rw, r)
		if err != nil {
			w.logger.Error(fmt.Sprintf("ERROR IN GOKIT WEB: %v", err))
		}
	}

	w.router.Handle(method, path, h)
}

// HandleAll registers handler on specified path and all of http methods
func (w *WebRouter) HandleAll(path string, handler Handler, middleware ...Middleware) {
	w.Handle("GET", path, handler, middleware...)
	w.Handle("HEAD", path, handler, middleware...)
	w.Handle("POST", path, handler, middleware...)
	w.Handle("PUT", path, handler, middleware...)
	w.Handle("DELETE", path, handler, middleware...)
	w.Handle("CONNECT", path, handler, middleware...)
	w.Handle("OPTIONS", path, handler, middleware...)
	w.Handle("TRACE", path, handler, middleware...)
	w.Handle("PATCH", path, handler, middleware...)
}

func (w *WebRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w.router.ServeHTTP(rw, r)
}
