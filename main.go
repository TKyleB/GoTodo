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

type AppConfig struct {
	usersHandler    users.UsersHandler
	snippetsHandler snippets.SnippetsHandler
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Error. DB_URL ENV is not set.")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error. PORT ENV is not set.")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database. %v", err)
	}
	err = db.Ping()
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
		Handler: corsMiddleware(mux),
		Addr:    ":" + port,
	}

	// Routes
	mux.HandleFunc("POST /api/users/register", appConfig.usersHandler.RegisterUser)
	mux.HandleFunc("POST /api/users/login", appConfig.usersHandler.LoginUser)
	mux.HandleFunc("POST /api/users/refresh", appConfig.usersHandler.RefreshUserToken)
	mux.HandleFunc("POST /api/users/logout", appConfig.usersHandler.LogoutUser)
	mux.HandleFunc("GET /api/users", appConfig.usersHandler.GetUser)

	mux.HandleFunc("POST /api/snippets", appConfig.snippetsHandler.CreateSnippet)
	mux.HandleFunc("GET /api/snippets", appConfig.snippetsHandler.GetSnippets)
	mux.HandleFunc("GET /api/snippets/{id}", appConfig.snippetsHandler.GetSnippetById)
	mux.HandleFunc("DELETE /api/snippets/{id}", appConfig.snippetsHandler.DeleteSnippetById)

	fmt.Printf("Starting server on %s\n", server.Addr)
	http.ListenAndServe(server.Addr, server.Handler)

}

func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "false")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}
