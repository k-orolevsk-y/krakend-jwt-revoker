package server

import "encoding/json"

type Response struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
}

func NewResponse(statusCode int, error string) []byte {
	resp := Response{StatusCode: statusCode, Error: error}

	jsonBytes, _ := json.Marshal(&resp)
	return jsonBytes
}
