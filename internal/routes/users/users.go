package users

import (
	"database/sql"
	"errors"
	"net/http"
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
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req CreateUserRequest
	err := utilites.DecodeJsonBody(w, r, &req)
	if err != nil {
		return
	}

	if !utilites.IsValidUsername(req.Username) {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "username must be >3 characters and not contain spaces or special characters")
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
	user, err := u.DbQueries.CreateUser(r.Context(), database.CreateUserParams{Username: req.Username, HashedPassword: hashedPassword})
	userResponse := auth.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Username: user.Username}
	if err != nil {
		// If error is non-unique username
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			utilites.ResponseWithError(w, r, http.StatusConflict, "username is already registered")
			return
		}
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
	}
	utilites.ResponseWithJson(w, r, http.StatusCreated, &userResponse)

}
func (u *UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	type LoginUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type LoginUserResponse struct {
		ID        uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Token     string    `json:"token"`
		Refresh   string    `json:"refresh_token"`
	}
	params := LoginUserRequest{}
	err := utilites.DecodeJsonBody(w, r, &params)
	if err != nil {
		return
	}

	// Get user if username exists
	user, err := u.DbQueries.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "invalid username/password combination")
		return
	}
	// Check if password matches stored
	err = u.AuthService.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "invalid username/password combination")
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
		Username:  user.Username,
		Token:     token,
		Refresh:   refreshToken.Token})

}
func (u *UsersHandler) RefreshUserToken(w http.ResponseWriter, r *http.Request) {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	type RefreshTokenResponse struct {
		Token string `json:"token"`
	}
	// Get refresh token from JSON request
	params := RefreshTokenRequest{}
	err := utilites.DecodeJsonBody(w, r, &params)
	if err != nil {
		return
	}

	// Check database against unexpired and revoked tokens
	refreshToken, err := u.DbQueries.GetRefreshToken(r.Context(), params.RefreshToken)
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
func (u *UsersHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	type LogoutUserRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	params := LogoutUserRequest{}
	err := utilites.DecodeJsonBody(w, r, &params)
	if err != nil {
		return
	}

	u.DbQueries.RevokeRefreshToken(r.Context(), params.RefreshToken)
	w.WriteHeader(http.StatusNoContent)

}

func (u *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := u.AuthService.GetAuthenticatedUser(r)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusOK, auth.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Username: user.Username})
}
