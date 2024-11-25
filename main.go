package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/cyberkillua/rssfeedagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	apiConfig := apiConfig{
		DB: database.New(connection),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiConfig.handlerCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUserByAPIKey))
	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.getFeeds)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server listening on port %v", portString)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
