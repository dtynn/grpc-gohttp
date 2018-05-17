package webapi

import "net/http"

func newResponseWritter(rw http.ResponseWriter) *responseWriter {
	if resw, ok := rw.(*responseWriter); ok {
		return resw
	}

	return &responseWriter{
		ResponseWriter: rw,
	}
}

type responseWriter struct {
	written bool
	http.ResponseWriter
}

func (r *responseWriter) Written() bool {
	return r.written
}

func (r *responseWriter) WriteHeader(statusCode int) {
	if r.written {
		return
	}

	r.written = true
	r.ResponseWriter.WriteHeader(statusCode)
}

// Written checks if a response writter has written
func Written(rw http.ResponseWriter) bool {
	type w interface {
		Written() bool
	}

	if ww, ok := rw.(w); ok {
		return ww.Written()
	}

	return false
}
