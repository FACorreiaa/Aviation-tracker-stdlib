package health

import (
	"encoding/json"
	"net/http"
)

type HandlerHealth struct{}

func (h HandlerHealth) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
