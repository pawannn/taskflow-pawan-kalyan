package engine

import (
	"encoding/json"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
)

type ErrorResponse struct {
	ReqID      string            `json:"req_id"`
	StatusCode int               `json:"status_code"`
	Error      string            `json:"error"`
	Message    string            `json:"message,omitempty"`
	Fields     map[string]string `json:"fields,omitempty"`
}

type Response struct {
	ReqID         string      `json:"req_id"`
	StatusCode    int         `json:"status_code"`
	ClientMessage string      `json:"client_message"`
	Data          interface{} `json:"data"`
}

func (e *HttpEngine) SendResponse(w http.ResponseWriter, reqID string, statusCode int, clientMessage string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		ReqID:         reqID,
		StatusCode:    statusCode,
		ClientMessage: clientMessage,
		Data:          data,
	}

	json.NewEncoder(w).Encode(response)
}

func (e *HttpEngine) SendErrorResponse(w http.ResponseWriter, reqID string, statusCode int, errorMessage string, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	status := ""

	switch statusCode {
	case http.StatusForbidden:
		status = domain.ErrForbidded
	case http.StatusUnauthorized:
		status = domain.ErrUnAuthorized
	case http.StatusNotFound:
		status = domain.ErrNotFound
	case http.StatusInternalServerError:
		status = domain.ErrInternalError
	}

	response := ErrorResponse{
		ReqID:      reqID,
		StatusCode: statusCode,
		Error:      status,
		Message:    errorMessage,
		Fields:     fields,
	}

	json.NewEncoder(w).Encode(response)
}
