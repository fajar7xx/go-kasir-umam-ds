package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	testData := map[string]string{"message": "success"}

	SendSuccess(w, testData, http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response SuccessResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, "success", data["message"])
}

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()

	SendError(w, "TEST_ERROR", "Test error message", http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "TEST_ERROR", response.Error.Code)
	assert.Equal(t, "Test error message", response.Error.Message)
	assert.Nil(t, response.Error.Details)
}

func TestSendErrorWithDetails(t *testing.T) {
	w := httptest.NewRecorder()
	details := map[string]string{"field": "name", "issue": "required"}

	SendErrorWithDetails(w, "VALIDATION_ERROR", "Validation failed", details, http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "VALIDATION_ERROR", response.Error.Code)
	assert.Equal(t, "Validation failed", response.Error.Message)
	assert.NotNil(t, response.Error.Details)

	detailsMap := response.Error.Details.(map[string]interface{})
	assert.Equal(t, "name", detailsMap["field"])
	assert.Equal(t, "required", detailsMap["issue"])
}

func TestParseIdFromPath_Success(t *testing.T) {
	mux := http.NewServeMux()
	var capturedID int

	mux.HandleFunc("GET /test/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := ParseIdFromPath(r, "id")
		assert.NoError(t, err)
		capturedID = id
	})

	req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, 123, capturedID)
}

func TestParseIdFromPath_InvalidFormat(t *testing.T) {
	mux := http.NewServeMux()
	var capturedError error

	mux.HandleFunc("GET /test/{id}", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParseIdFromPath(r, "id")
		capturedError = err
	})

	req := httptest.NewRequest(http.MethodGet, "/test/abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Error(t, capturedError)
}

func TestParseIdFromPath_EmptyPath(t *testing.T) {
	mux := http.NewServeMux()
	var capturedError error

	mux.HandleFunc("GET /test/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Try to parse a non-existent parameter name, which returns empty string
		_, err := ParseIdFromPath(r, "nonexistent")
		capturedError = err
	})

	req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Error(t, capturedError)
}
