package snippets

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TKyleB/GoTodo/internal/auth"
	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/utilites"
	"github.com/google/uuid"
)

type SnippetsHandler struct {
	DbQueries   *database.Queries
	TokenSecret string
}
type Snippet struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Language    string    `json:"language"`
	AuthorID    uuid.UUID `json:"author_id"`
	SnippetText string    `json:"text"`
}

func (s *SnippetsHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Language string `json:"language"`
		Text     string `json:"text"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "auth token not provided")
		return
	}
	userID, err := auth.ValidateJWT(token, s.TokenSecret)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "invalid token")
		return
	}

	params := parameters{}
	err = utilites.DecodeJsonBody(w, r, &params)
	if err != nil {
		return
	}
	if len(params.Text) == 0 {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "text body is empty")
		return
	}
	languageID, err := s.DbQueries.GetLanguageByName(r.Context(), params.Language)
	if err != nil {
		errorText := fmt.Sprintf("language: %s is not currently supported", params.Language)
		utilites.ResponseWithError(w, r, http.StatusBadRequest, errorText)
		return
	}
	snippet, err := s.DbQueries.CreateSnippet(r.Context(), database.CreateSnippetParams{LanguageID: languageID, AuthorID: userID, SnippetText: params.Text})
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusCreated, Snippet{ID: snippet.ID, CreatedAt: snippet.CreatedAt, UpdatedAt: snippet.UpdatedAt, Language: params.Language, AuthorID: snippet.AuthorID, SnippetText: snippet.SnippetText})

}
