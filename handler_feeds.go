package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cyberkillua/rssfeedagg/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error creating feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiConfig *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error getting feeds: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}


