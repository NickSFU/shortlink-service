package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSONError(
	w http.ResponseWriter,
	message string,
	status int,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(
		ErrorResponse{
			Error: message,
		},
	)
}
