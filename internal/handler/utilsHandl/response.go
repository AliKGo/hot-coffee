package utilsHandl

import (
	"encoding/json"
	"net/http"
)

func SendJSONMessages(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"message": message,
		"status":  statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "JSON serialization error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
