package context

import (
	"context"
	"net/http"
)

var (
	requestKey  struct{}
	responseKey struct{}
)

// NewWithRequest returns context.Context with *http.Request
func NewWithRequest(parent context.Context, req *http.Request) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	return context.WithValue(parent, requestKey, req)
}

// NewWithResponseWriter returns context.Context with http.ResponseWriter
func NewWithResponseWriter(parent context.Context, rw http.ResponseWriter) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	return context.WithValue(parent, responseKey, rw)
}

// RequestFromCtx returns request in context
func RequestFromCtx(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(requestKey).(*http.Request)
	return req, ok
}

// ResponseWriterFromCtx returns response writer in context
func ResponseWriterFromCtx(ctx context.Context) (http.ResponseWriter, bool) {
	rw, ok := ctx.Value(responseKey).(http.ResponseWriter)
	return rw, ok
}
