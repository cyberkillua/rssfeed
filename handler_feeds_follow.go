package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cyberkillua/rssfeedagg/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	feedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error creating feed follow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}
