package snippets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TKyleB/GoTodo/internal/auth"
	"github.com/TKyleB/GoTodo/internal/database"
	"github.com/TKyleB/GoTodo/internal/utilites"
	"github.com/google/uuid"
)

type SnippetsHandler struct {
	DbQueries   *database.Queries
	AuthService *auth.AuthService
}
type Snippet struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Language    string    `json:"language"`
	UserID      uuid.UUID `json:"author_id"`
	SnippetText string    `json:"text"`
}

type SnippetsResponse struct {
	Count    int32
	Next     *string
	Previous *string
	Results  []Snippet
}

func (s *SnippetsHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Language string `json:"language"`
		Text     string `json:"text"`
	}

	user, err := s.AuthService.GetAuthenticatedUser(r)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, err.Error())
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
	snippet, err := s.DbQueries.CreateSnippet(r.Context(), database.CreateSnippetParams{LanguageID: languageID, UserID: user.ID, SnippetText: params.Text})
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusCreated, Snippet{ID: snippet.ID, CreatedAt: snippet.CreatedAt, UpdatedAt: snippet.UpdatedAt, Language: params.Language, UserID: snippet.UserID, SnippetText: snippet.SnippetText})

}
func (s *SnippetsHandler) GetSnippets(w http.ResponseWriter, r *http.Request) {
	limit := int32(10)
	offset := int32(0)
	var language sql.NullString
	languageString := r.URL.Query().Get("language")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	if limitString != "" {
		if parseLimit, err := strconv.ParseInt(limitString, 10, 32); err == nil {
			limit = int32(parseLimit)
		}
	}
	if offsetString != "" {
		if parseOffset, err := strconv.ParseInt(offsetString, 10, 32); err == nil {
			offset = int32(parseOffset)
		}
	}
	if languageString != "" {
		language.Scan(languageString)
	}

	dbSnippets, _ := s.DbQueries.GetSnippetsByCreatedAt(r.Context(), database.GetSnippetsByCreatedAtParams{Limit: limit, Offset: offset, Language: language})
	snippetsCount, _ := s.DbQueries.GetSnippetCount(r.Context())
	count := int32(snippetsCount)

	var snippets []Snippet
	for _, snippet := range dbSnippets {
		snippets = append(snippets, Snippet{ID: snippet.ID, CreatedAt: snippet.CreatedAt, UpdatedAt: snippet.UpdatedAt, Language: snippet.Language, UserID: snippet.UserID, SnippetText: snippet.SnippetText})
	}
	baseURL := fmt.Sprintf("%s://%s", r.URL.Scheme, r.Host)
	var next *string
	var previous *string

	if !(offset+limit >= count) {
		nextURL := fmt.Sprintf("%s/api/snippets?limit=%v&offset=%v", baseURL, limit, offset+limit)
		next = &nextURL
	}
	if offset > 0 {
		prevURL := fmt.Sprintf("%s/api/snippets?limit=%v&offset=%v", baseURL, limit, max(offset-limit, 0))
		previous = &prevURL
	}
	response := SnippetsResponse{Count: count, Next: next, Previous: previous, Results: snippets}

	utilites.ResponseWithJson(w, r, http.StatusOK, &response)
}
