package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessageResponse struct {
	Message string `json:"message"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("WriteJSONResponse error to encode data: %s", err.Error())
	}
}

func WriteJSONResponseError(w http.ResponseWriter, status int, error string) {
	WriteJSONResponse(w, status, ErrorMessageResponse{Message: error})
}

func DecodeJSON(r *http.Request, value any) error {
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	if err := decode.Decode(&value); err != nil {
		log.Printf("decode json error: %v, value: %v", err, value)
		return err
	}

	return nil
}
