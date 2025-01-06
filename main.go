package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/routes/users"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Used to connect to DB
)

const PORT = "8080"

type AppConfig struct {
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
	appConfig := AppConfig{
		usersHandler: users.UsersHandler{DbQueries: dbQueries, TokenSecret: tokenSecret},
	}
	server := http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	// Routes
	mux.HandleFunc("POST /api/users/register", appConfig.usersHandler.RegisterUser)
	mux.HandleFunc("POST /api/users/login", appConfig.usersHandler.LoginUser)
	mux.HandleFunc("POST /api/users/refresh", appConfig.usersHandler.RefreshUserToken)
	mux.HandleFunc("GET /api/users/", appConfig.usersHandler.GetUser)

	fmt.Printf("Starting server on %s\n", server.Addr)
	server.ListenAndServe()

}
