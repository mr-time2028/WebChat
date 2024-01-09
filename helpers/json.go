package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/mr-time2028/WebChat/validators"
)

// ReadJSON read request and extract request payload from it
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) validators.Validation {
	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// validation
	validator := validators.New()
	validator.JsonValidation(r, data)
	return validator
}

// WriteJSON write data to output
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorMapJSON write an error map message as a json to the output
func ErrorMapJSON(w http.ResponseWriter, error validators.Errors) {
	// {
	//		"error": true
	// 		"message": {
	//			"email": [
	//				"first error for email",
	//				"second error for email",
	//			],
	//			"password": [
	//				"some error for password"
	//			]
	//      }
	// }
	type jsonResponse struct {
		Error   bool                `json:"error"`
		Message map[string][]string `json:"message"`
	}

	var payload jsonResponse
	var statusCode int

	payload.Error = true
	payload.Message = error.MessageMap
	statusCode = error.Code

	if statusCode == http.StatusInternalServerError {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err := WriteJSON(w, statusCode, payload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// ErrorStrJSON write a str error message as a json to the output
func ErrorStrJSON(w http.ResponseWriter, error error, status ...int) {
	// {
	//		"error": true,
	//    	"message": "some error"
	// }
	type jsonResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = error.Error()

	err := WriteJSON(w, statusCode, payload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
