package handle

import (
	"bytes"
	"encoding/json"
)

func ApiResponse(statusCode int, body interface{}) (Response, error) {
	var buf bytes.Buffer

	resp := Response{Headers: map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Methods": "GET, POST"}}
	resp.StatusCode = statusCode

	respBody, _ := json.Marshal(body)

	json.HTMLEscape(&buf, respBody)
	resp.Body = string(respBody)

	return resp, nil
}
