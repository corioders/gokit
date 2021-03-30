package web

import (
	"fmt"
	"net/http"

	"github.com/corioders/gokit/log"
	"github.com/dimfeld/httptreemux"
)

type Router interface {
	http.Handler
	Handle(method string, path string, handler Handler, middleware ...Middleware)
	HandleAll(path string, handler Handler, middleware ...Middleware)
	NewGroup(path string, middleware ...Middleware) RouterGroup
}

type RouterGroup interface {
	Handle(method string, path string, handler Handler, middleware ...Middleware)
	HandleAll(path string, handler Handler, middleware ...Middleware)
	NewGroup(path string, middleware ...Middleware) RouterGroup
}

type internalRouter struct {
	logger log.Logger

	router     *httptreemux.ContextMux
	middleware []Middleware
}

func NewRouter(logger log.Logger, middleware ...Middleware) Router {
	return &internalRouter{
		logger: logger,

		router:     httptreemux.NewContextMux(),
		middleware: middleware,
	}
}

// Handle registers handler on specified path and method
func (ir *internalRouter) Handle(method string, path string, handler Handler, middleware ...Middleware) {
	handle(ir.router, ir.logger, method, path, handler, ir.middleware, middleware)
}

// HandleAll registers handler on specified path and all of http methods
func (ir *internalRouter) HandleAll(path string, handler Handler, middleware ...Middleware) {
	handleAll(ir, path, handler, middleware)
}

func (wr *internalRouter) NewGroup(path string, middleware ...Middleware) RouterGroup {
	group := wr.router.NewContextGroup(path)
	return newInternalGroup(wr.logger, path, group, middleware...)
}

func (ir *internalRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ir.router.ServeHTTP(rw, r)
}

type internalGroup struct {
	logger log.Logger

	group      *httptreemux.ContextGroup
	middleware []Middleware
}

func newInternalGroup(logger log.Logger, path string, group *httptreemux.ContextGroup, middleware ...Middleware) *internalGroup {
	return &internalGroup{
		logger: logger,

		group:      group,
		middleware: middleware,
	}
}

// Handle registers handler on specified path and method
func (ig *internalGroup) Handle(method string, path string, handler Handler, middleware ...Middleware) {
	handle(ig.group, ig.logger, method, path, handler, ig.middleware, middleware)
}

// HandleAll registers handler on specified path and all of http methods
func (ig *internalGroup) HandleAll(path string, handler Handler, middleware ...Middleware) {
	handleAll(ig, path, handler, middleware)
}

func (ig *internalGroup) NewGroup(path string, middleware ...Middleware) RouterGroup {
	group := ig.group.NewContextGroup(path)
	return newInternalGroup(ig.logger, path, group, middleware...)
}

type handlerAsHandlerFunc interface {
	Handle(method, path string, handler http.HandlerFunc)
}

func handle(object handlerAsHandlerFunc, logger log.Logger, method string, path string, handler Handler, generalMiddleware, specificMiddleware []Middleware) {
	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(specificMiddleware, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(generalMiddleware, handler)

	h := func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := handler(ctx, rw, r)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR IN GOKIT WEB: %v", err))
		}
	}

	object.Handle(method, path, h)
}

type handlerAsWebHandler interface {
	Handle(method, path string, handler Handler, middleware ...Middleware)
}

func handleAll(object handlerAsWebHandler, path string, handler Handler, middleware []Middleware) {
	object.Handle("GET", path, handler, middleware...)
	object.Handle("HEAD", path, handler, middleware...)
	object.Handle("POST", path, handler, middleware...)
	object.Handle("PUT", path, handler, middleware...)
	object.Handle("DELETE", path, handler, middleware...)
	object.Handle("CONNECT", path, handler, middleware...)
	object.Handle("OPTIONS", path, handler, middleware...)
	object.Handle("TRACE", path, handler, middleware...)
	object.Handle("PATCH", path, handler, middleware...)
}
