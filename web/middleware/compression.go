package middleware

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/web"
)

type compressor interface {
	io.WriteCloser
	Reset(io.Writer)
}

type compressionResponseWriter struct {
	http.ResponseWriter

	compressionWriter compressor
	statusCode        int
	headerWritten     bool
}

var gzipPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		return w
	},
}

var flatePool = sync.Pool{
	New: func() interface{} {
		w, _ := flate.NewWriter(nil, flate.BestSpeed)
		return w
	},
}

func (crw *compressionResponseWriter) WriteHeader(statusCode int) {
	crw.statusCode = statusCode
	crw.headerWritten = true

	crw.ResponseWriter.WriteHeader(statusCode)
}

func (crw *compressionResponseWriter) Write(b []byte) (int, error) {
	if !crw.headerWritten {
		// This is exactly what Go would also do if it hasn't been written yet.
		crw.WriteHeader(http.StatusOK)
	}
	return crw.compressionWriter.Write(b)
}

func Compression() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
			if r.Header.Get("upgrade") != "" {
				return handler(ctx, rw, r)
			}

			usingCompression := false
			var compressionWriter compressor
			compressions := strings.Split(r.Header.Get("Accept-Encoding"), ",")
			for _, compression := range compressions {
				switch strings.Trim(compression, " ") {
				case "*":
				case "gzip":
					usingCompression = true
					rw.Header().Set("Content-Encoding", "gzip")
					gzipWriter := gzipPool.Get().(*gzip.Writer)
					gzipWriter.Reset(rw)
					compressionWriter = gzipWriter

				case "deflate":
					usingCompression = true
					rw.Header().Set("Content-Encoding", "deflate")
					flateWriter := flatePool.Get().(*flate.Writer)
					flateWriter.Reset(rw)
					compressionWriter = flateWriter
				}

				// We want to apply only one compression.
				if usingCompression {
					break
				}
			}

			if !usingCompression {
				// Just do nothing.
				return handler(ctx, rw, r)
			}

			crw := compressionResponseWriter{compressionWriter: compressionWriter, ResponseWriter: rw}
			err := handler(ctx, &crw, r)
			if err != nil {
				crw.compressionWriter.Reset(nil)
				// Don't mess with error returned by handler.
				return err
			}

			if crw.statusCode != http.StatusNotModified && crw.statusCode != http.StatusNoContent {
				if err := crw.compressionWriter.Close(); err != nil {
					return errors.WithMessage(err, "closing compressor")
				}
			} else {
				// No data must be send.
				crw.compressionWriter.Reset(nil)
			}

			switch w := crw.compressionWriter.(type) {
			case *gzip.Writer:
				gzipPool.Put(w)

			case *flate.Writer:
				flatePool.Put(w)
			}

			return nil
		}
	}
}
