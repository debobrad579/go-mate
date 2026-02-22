package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/debobrad579/chessgo/internal/auth"
)

func (cfg *Config) ApiMeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("null"))
		return
	}

	userID, err := auth.ValidateJWT(cookie.Value, cfg.TokenSecret)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
