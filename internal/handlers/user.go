package handlers

import (
	"net/http"

	"github.com/debobrad579/chessgo/internal/auth"
	"github.com/debobrad579/chessgo/internal/database"
)

func (cfg *Config) getUser(r *http.Request) (*database.User, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, nil
	}

	userID, err := auth.ValidateJWT(cookie.Value, cfg.TokenSecret)
	if err != nil {
		return nil, err
	}

	user, err := cfg.DB.GetUser(r.Context(), userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
