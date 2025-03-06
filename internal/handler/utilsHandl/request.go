package utilsHandl

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseJSONBody(r *http.Request, v interface{}) (string, int) {
	if r.Body == nil {
		return "Empty request body", http.StatusBadRequest
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "Failed to read request body", http.StatusInternalServerError
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return "Invalid JSON format", http.StatusBadRequest
	}

	return "OK", http.StatusOK
}
