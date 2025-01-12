package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TKyleB/GoTodo/internal/auth"
	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/routes/snippets"
	"github.com/TKyleB/GoTodo/internal/routes/users"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Used to connect to DB
)

const PORT = "8080"

type AppConfig struct {
	usersHandler    users.UsersHandler
	snippetsHandler snippets.SnippetsHandler
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database. %v", err)
	}
	// Services
	dbQueries := database.New(db)
	authService := auth.AuthService{
		TokenSecret:                os.Getenv("TOKEN_SECRET"),
		TokenExpirationTime:        time.Minute * 10,
		RefreshTokenExpirationTime: time.Hour * 24 * 30,
		Issuer:                     "snippetz",
		DbQueries:                  dbQueries,
	}

	appConfig := AppConfig{
		usersHandler:    users.UsersHandler{DbQueries: dbQueries, AuthService: &authService},
		snippetsHandler: snippets.SnippetsHandler{DbQueries: dbQueries, AuthService: &authService},
	}

	server := http.Server{
		Handler: mux,
		Addr:    ":" + PORT,
	}

	// Routes
	mux.HandleFunc("POST /api/users/register", appConfig.usersHandler.RegisterUser)
	mux.HandleFunc("POST /api/users/login", appConfig.usersHandler.LoginUser)
	mux.HandleFunc("POST /api/users/refresh", appConfig.usersHandler.RefreshUserToken)
	mux.HandleFunc("GET /api/users", appConfig.usersHandler.GetUser)

	mux.HandleFunc("POST /api/snippets", appConfig.snippetsHandler.CreateSnippet)
	mux.HandleFunc("GET /api/snippets", appConfig.snippetsHandler.GetSnippets)

	fmt.Printf("Starting server on %s\n", server.Addr)
	http.ListenAndServe(server.Addr, corsMiddleware(mux))

}

func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}
