package handlers

import (
	"encoding/json"
	"net/http"
)

func (cfg *Config) ApiMeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.getUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("null"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
