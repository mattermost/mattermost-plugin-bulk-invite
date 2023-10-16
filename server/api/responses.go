package api

import (
	"fmt"
	"net/http"
)

type responseOption func(w http.ResponseWriter)

func sendResponse(w http.ResponseWriter, opts ...responseOption) {
	for _, opt := range opts {
		opt(w)
	}
}

func withStatusCode(statusCode int) responseOption {
	return func(w http.ResponseWriter) {
		w.WriteHeader(statusCode)
	}
}

func withHeader(key, value string) responseOption {
	return func(w http.ResponseWriter) {
		w.Header().Set(key, value)
	}
}

func withBody(body string, args ...any) responseOption {
	return func(w http.ResponseWriter) {
		_, _ = w.Write([]byte(fmt.Sprintf(body, args...)))
	}
}
