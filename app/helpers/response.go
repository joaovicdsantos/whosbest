package helpers

import (
	"encoding/json"
	"net/http"
	"time"
)

type DataType interface{}

type Response struct {
	Data       DataType
	StatusCode int
	Timestamp  int64
}

func NewResponse(data DataType, statusCode int) Response {
	var response Response

	response.Data = data
	response.StatusCode = statusCode
	response.Timestamp = time.Now().Unix()

	return response
}

func (r *Response) SendResponse(w http.ResponseWriter) {
	response, err := json.Marshal(map[string]interface{}{
		"data":        r.Data,
		"status_code": r.StatusCode,
		"timestamp":   r.Timestamp,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	w.Write(response)
}
