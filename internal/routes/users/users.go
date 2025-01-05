package users

import (
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/utilites"
	"github.com/google/uuid"
)

type UsersHandler struct {
	DbQueries *database.Queries
}
type User struct {
	ID        uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	user, err := u.DbQueries.CreateUser(r.Context(), req.Email)
	userResponse := User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		fmt.Printf("%v", err)
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusCreated, &userResponse)
}
