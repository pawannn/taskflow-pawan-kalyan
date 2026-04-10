package engine

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	ReqID         string      `json:"req_id"`
	StatusCode    int         `json:"status_cdoe"`
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
