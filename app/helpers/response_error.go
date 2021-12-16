package helpers

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorType interface{}

type ResponseError struct {
	Error      ErrorType
	StatusCode int
	Timestamp  int64
}

func NewResponseError(err ErrorType, statusCode int) ResponseError {
	var response ResponseError

	response.Error = err
	response.StatusCode = statusCode
	response.Timestamp = time.Now().Unix()

	return response
}

func (r *ResponseError) SendResponse(w http.ResponseWriter) {
	response, err := json.Marshal(map[string]interface{}{
		"error":       r.Error,
		"status_code": r.StatusCode,
		"timestamp":   r.Timestamp,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(r.StatusCode)
	w.Write(response)
}
