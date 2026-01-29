package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

func SendSuccess(w http.ResponseWriter, data interface{}, status int) {
	response := SuccessResponse{Data: data}
	sendJSON(w, status, response)
}

func SendError(w http.ResponseWriter, code, message string, status int) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	sendJSON(w, status, response)
}

func SendErrorWithDetails(w http.ResponseWriter, code, message string, details interface{}, status int) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	sendJSON(w, status, response)
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ParseIdFromPath(r *http.Request, paramName string) (int, error) {
	idStr := r.PathValue(paramName)
	return strconv.Atoi(idStr)
}
