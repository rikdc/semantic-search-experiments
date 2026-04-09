package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandleRequest is the top-level HTTP handler. It routes the request to the
// appropriate subhandler based on the URL path and HTTP method.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case "/api/v1/items":
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		WriteJSON(w, http.StatusOK, map[string]any{"items": []string{}})
	default:
		WriteError(w, http.StatusNotFound, fmt.Sprintf("no route for %s", r.URL.Path))
	}
}

// WriteJSON serialises v as JSON and writes it to w with the given status code.
// Sets Content-Type to application/json.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encoding response", http.StatusInternalServerError)
	}
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
