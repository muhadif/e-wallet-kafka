package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func ResponseSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	resp := Response{
		Status: "success",
		Data:   data,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
	return
}

type ErrorResponse struct {
	Status    string `json:"status,omitempty"`
	ErrorCode int    `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
}

func ResponseError(w http.ResponseWriter, statusCode int, err error) {
	resp := &ErrorResponse{
		Status:    "error",
		ErrorCode: statusCode,
		Message:   err.Error(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
	return
}
