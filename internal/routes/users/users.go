package users

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/TKyleB/GoTodo/internal/auth"
	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/utilites"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UsersHandler struct {
	DbQueries   *database.Queries
	AuthService *auth.AuthService
}

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
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "invalid email format")
		return
	}
	if len(req.Password) < 6 {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "password must be 6 or greater characters")
		return
	}
	hashedPassword, err := u.AuthService.HashPassword(req.Password)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "invalid characters in password")
		return
	}
	user, err := u.DbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: req.Email, HashedPassword: hashedPassword})
	userResponse := auth.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	if err != nil {
		// If error is non-unique email
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			utilites.ResponseWithError(w, r, http.StatusConflict, "email is already registered")
			return
		}
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
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
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "invalid email/password combination")
		return
	}
	// Check if password matches stored
	err = u.AuthService.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "invalid email/password combination")
		return
	}
	// Create JWT token
	token, err := u.AuthService.MakeJWT(user.ID)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
		return
	}
	// Create Refresh Token
	refreshTokenString, err := u.AuthService.MakeRefreshToken()
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
		return
	}
	refreshToken, err := u.DbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(u.AuthService.RefreshTokenExpirationTime),
		RevokedAt: sql.NullTime{}})
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
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
func (u *UsersHandler) RefreshUserToken(w http.ResponseWriter, r *http.Request) {
	type RefreshTokenResponse struct {
		Token string `json:"token"`
	}

	// Get refresh token from auth headers
	refreshTokenString, err := u.AuthService.GetBearerToken(r.Header)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "invalid auth headers")
		return
	}
	// Check database against unexpired and revoked tokens
	refreshToken, err := u.DbQueries.GetRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "expired or invalid token")
		return
	}

	// Generate new JWT token
	newToken, err := u.AuthService.MakeJWT(refreshToken.UserID)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
		return
	}
	refreshTokenResponse := RefreshTokenResponse{Token: newToken}
	utilites.ResponseWithJson(w, r, http.StatusOK, &refreshTokenResponse)
}
func (u *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := u.AuthService.GetAuthenticatedUser(r)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusOK, auth.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
}
