package web

import (
	"context"
	"net/http"
)

type Handler func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error
