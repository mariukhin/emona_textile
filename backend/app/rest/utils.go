package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type restResponseErrorPayload struct {
	Errors []string `json:"error"`
}

func httpJsonResponse(w http.ResponseWriter, resp interface{}) {
	httpJsonResponseCode(w, resp, http.StatusOK)
}

func httpJsonResponseCode(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"error": [
				"internal"
			]
		}`)
		return
	}
}

func httpResponseErrors(w http.ResponseWriter, statusCode int, errors []error) {
	errs := make([]string, len(errors))
	for idx, e := range errors {
		errs[idx] = e.Error()
	}

	resp := restResponseErrorPayload{
		Errors: errs,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"error": [
				"internal"
			]
		}`)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	w.Write(respBytes)
}

func httpResponseError(w http.ResponseWriter, statusCode int, error string) {
	resp := restResponseErrorPayload{
		Errors: []string{error},
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"error": [
				"internal"
			]
		}`)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	w.Write(respBytes)
}

func httpPlainError(w http.ResponseWriter, statusCode int, error string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(error))
}

func tokenFromHeader(r *http.Request) string {
	// Get token from X-Auth-Token header.
	token := r.Header.Get("X-Auth-Token")
	return token
}
