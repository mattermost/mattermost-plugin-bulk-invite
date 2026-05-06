package api

import (
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

func withBody(body string) responseOption {
	return func(w http.ResponseWriter) {
		_, _ = w.Write([]byte(body))
	}
}

func sendInternalServerError(w http.ResponseWriter) {
	sendResponse(w, withStatusCode(http.StatusInternalServerError), withBody(`{"error": "internal server error"}`))
}
