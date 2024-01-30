package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	sendInternalServerError(w)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestSendResponse(t *testing.T) {
	expectedContentType := "application/json"
	expectedBody := `{"message": "success"}`
	expectedStatusCode := http.StatusOK

	w := httptest.NewRecorder()
	sendResponse(w, withStatusCode(expectedStatusCode), withHeader("Content-Type", expectedContentType), withBody(expectedBody))

	require.Equalf(t, expectedStatusCode, w.Code, "Expected status code %d, but got %d", expectedStatusCode, w.Code)

	actualContentType := w.Header().Get("Content-Type")
	require.Equalf(t, expectedContentType, actualContentType, "Expected Content-Type header %q, but got %q", expectedContentType, actualContentType)

	actualBody := w.Body.String()
	require.Equalf(t, expectedBody, actualBody, "Expected response body %q, but got %q", expectedBody, actualBody)
}
