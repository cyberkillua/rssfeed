package main

import (
	"fmt"
	"net/http"

	"github.com/cyberkillua/rssfeedagg/internal/auth"
	"github.com/cyberkillua/rssfeedagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apikey)

		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprint("Error getting user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
