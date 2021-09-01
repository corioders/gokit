package web

import (
	"fmt"
	"net/http"

	"github.com/corioders/gokit/log"
	"github.com/dimfeld/httptreemux"
)

type RouterGroup interface {
	Handle(method string, path string, handler Handler, middleware ...Middleware)
	HandleAll(path string, handler Handler, middleware ...Middleware)
	NewGroup(path string, middleware ...Middleware) RouterGroup
}
type Router interface {
	RouterGroup
	http.Handler
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
	handleAll(ir.router, ir.logger, path, handler, ir.middleware, middleware)
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
	handleAll(ig.group, ig.logger, path, handler, ig.middleware, middleware)
}

func (ig *internalGroup) NewGroup(path string, middleware ...Middleware) RouterGroup {
	group := ig.group.NewContextGroup(path)
	return newInternalGroup(ig.logger, path, group, middleware...)
}

type router interface {
	Handle(method, path string, handler http.HandlerFunc)
}

func handle(r router, logger log.Logger, method string, path string, handler Handler, generalMiddleware, specificMiddleware []Middleware) {
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

	r.Handle(method, path, h)
}

func handleAll(r router, logger log.Logger, path string, handler Handler, generalMiddleware, specificMiddleware []Middleware) {
	handle(r, logger, http.MethodGet, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodHead, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodPost, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodPut, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodPatch, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodDelete, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodConnect, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodOptions, path, handler, generalMiddleware, specificMiddleware)
	handle(r, logger, http.MethodTrace, path, handler, generalMiddleware, specificMiddleware)
}
