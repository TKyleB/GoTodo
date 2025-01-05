package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/TKyleB/GoTodo/internal/auth"
	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/utilites"
	"github.com/google/uuid"
)

type UsersHandler struct {
	DbQueries   *database.Queries
	TokenSecret string
}
type User struct {
	ID        uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

const TOKEN_EXPIRATION_TIME = time.Minute * 10            // 10 Minutes
const REFRESH_TOKEN_EXPIRATION_TIME = time.Hour * 24 * 30 // 30 Days

func (u *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	type CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req CreateUserRequest
	err := utilites.DecodeJsonBody(w, r, &req)
	if err != nil {
		return
	}
	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be 6 or greater characters long", http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error with password", http.StatusBadRequest)
		return
	}
	user, err := u.DbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: req.Email, HashedPassword: hashedPassword})
	userResponse := User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		fmt.Printf("%v", err)
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusCreated, &userResponse)

}
func (u *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	type LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type LoginUserResponse struct {
		ID        uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
		Refresh   string    `json:"refresh_token"`
	}
	params := LoginUserRequest{}
	err := utilites.DecodeJsonBody(w, r, &params)
	if err != nil {
		return
	}

	// Get user if email exists
	user, err := u.DbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		http.Error(w, "Invalid username/password", http.StatusUnauthorized)
		return
	}
	// Check if password matches stored
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		http.Error(w, "Invalid username/password", http.StatusUnauthorized)
		return
	}
	// Create JWT token
	token, err := auth.MakeJWT(user.ID, u.TokenSecret, TOKEN_EXPIRATION_TIME)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	// Create Refresh Token
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	refreshToken, err := u.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(REFRESH_TOKEN_EXPIRATION_TIME),
		RevokedAt: sql.NullTime{}})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusOK, LoginUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
		Refresh:   refreshToken.Token})

}
