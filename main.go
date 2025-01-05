package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/routes/users"
	_ "github.com/lib/pq" // Used to connect to DB
)

const PORT = "8080"

type ApiConfig struct {
	usersHandler users.UsersHandler
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database. %v", err)
	}
	dbQueries := database.New(db)
	apiConfig := ApiConfig{
		usersHandler: users.UsersHandler{DbQueries: dbQueries, TokenSecret: tokenSecret},
	}
	server := http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	// Routes
	mux.HandleFunc("POST /api/users/register", apiConfig.usersHandler.RegisterUser)
	mux.HandleFunc("POST /api/users/login", apiConfig.usersHandler.LoginUser)
	mux.HandleFunc("POST /api/users/refresh", apiConfig.usersHandler.RefreshUserToken)

	fmt.Printf("Starting server on %s\n", server.Addr)
	server.ListenAndServe()

}
