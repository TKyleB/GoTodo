package snippets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TKyleB/snippetz/internal/auth"
	"github.com/TKyleB/snippetz/internal/database"
	"github.com/TKyleB/snippetz/internal/utilites"
	"github.com/google/uuid"
)

type SnippetsHandler struct {
	DbQueries   *database.Queries
	AuthService *auth.AuthService
}
type Snippet struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Language     string    `json:"language"`
	UserID       uuid.UUID `json:"author_id"`
	UserName     string    `json:"username"`
	SnippetText  string    `json:"snippet_text"`
	SnippetDesc  string    `json:"snippet_desc"`
	SnippetTitle string    `json:"snippet_title"`
}

type Results struct {
	Snippets  []Snippet      `json:"snippets"`
	Languages map[string]int `json:"languages"`
}

type SnippetsResponse struct {
	Count    int32   `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  Results `json:"results"`
}

func (s *SnippetsHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Language     string `json:"language"`
		SnippetText  string `json:"snippet_text"`
		SnippetDesc  string `json:"Snippet_desc"`
		SnippetTitle string `json:"snippet_title"`
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
	if len(params.SnippetText) == 0 {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "snippet_text is empty")
		return
	}
	if len(params.SnippetDesc) == 0 {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "snippet_desc is empty")
		return
	}
	if len(params.SnippetTitle) == 0 {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "snippet_title is empty")
		return
	}
	languageID, err := s.DbQueries.GetLanguageByName(r.Context(), params.Language)
	if err != nil {
		errorText := fmt.Sprintf("language: %s is not currently supported", params.Language)
		utilites.ResponseWithError(w, r, http.StatusBadRequest, errorText)
		return
	}
	snippet, err := s.DbQueries.CreateSnippet(r.Context(), database.CreateSnippetParams{
		LanguageID:         languageID,
		UserID:             user.ID,
		SnippetText:        params.SnippetText,
		SnippetDescription: params.SnippetDesc,
		SnippetTitle:       params.SnippetTitle,
	})
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusInternalServerError, "server error")
		return
	}
	utilites.ResponseWithJson(w, r, http.StatusCreated, Snippet{
		ID:           snippet.ID,
		CreatedAt:    snippet.CreatedAt,
		UpdatedAt:    snippet.UpdatedAt,
		Language:     params.Language,
		UserID:       snippet.UserID,
		SnippetText:  snippet.SnippetText,
		SnippetDesc:  snippet.SnippetDescription,
		SnippetTitle: snippet.SnippetTitle,
		UserName:     snippet.Username,
	})

}
func (s *SnippetsHandler) GetSnippets(w http.ResponseWriter, r *http.Request) {
	// Set-up for pagination
	var count int32
	limit := int32(5)
	offset := int32(0)

	var language sql.NullString
	var username sql.NullString
	var search sql.NullString
	languageString := r.URL.Query().Get("language")
	usernameString := r.URL.Query().Get("username")
	searchString := r.URL.Query().Get("q")
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
	if usernameString != "" {
		username.Scan(usernameString)
	}
	if searchString != "" {
		words := strings.Join(strings.Split(searchString, " "), " & ")
		search.Scan(words)
	}

	dbSnippets, _ := s.DbQueries.GetSnippetsByCreatedAt(r.Context(), database.GetSnippetsByCreatedAtParams{Limit: limit, Offset: offset, Language: language, Username: username, Search: search})

	languageCounts := make(map[string]int)
	var snippets []Snippet
	for i, snippet := range dbSnippets {
		// If snippets. Update total count
		if i == 0 {
			count = int32(snippet.TotalCount)
		}

		languageCounts[snippet.Language] += 1

		snippets = append(snippets, Snippet{
			ID:           snippet.ID,
			CreatedAt:    snippet.CreatedAt,
			UpdatedAt:    snippet.UpdatedAt,
			Language:     snippet.Language,
			UserID:       snippet.UserID,
			SnippetText:  snippet.SnippetText,
			UserName:     snippet.Username,
			SnippetDesc:  snippet.SnippetDescription,
			SnippetTitle: snippet.SnippetTitle,
		})
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
	var results = Results{}
	results.Snippets = snippets
	results.Languages = languageCounts
	response := SnippetsResponse{Count: count, Next: next, Previous: previous, Results: results}

	utilites.ResponseWithJson(w, r, http.StatusOK, &response)
}

func (s *SnippetsHandler) GetSnippetById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "invalid id")
	}
	dbSnippet, _ := s.DbQueries.GetSnippetById(r.Context(), id)
	snippet := Snippet{
		ID:           dbSnippet.ID,
		CreatedAt:    dbSnippet.CreatedAt,
		UpdatedAt:    dbSnippet.UpdatedAt,
		UserName:     dbSnippet.Username,
		UserID:       dbSnippet.UserID,
		SnippetDesc:  dbSnippet.SnippetDescription,
		SnippetTitle: dbSnippet.SnippetTitle,
		SnippetText:  dbSnippet.SnippetText,
		Language:     dbSnippet.Language,
	}
	utilites.ResponseWithJson(w, r, http.StatusOK, &snippet)
}

func (s *SnippetsHandler) DeleteSnippetById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusBadRequest, "invalid id")
	}
	user, err := s.AuthService.GetAuthenticatedUser(r)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	snippet, err := s.DbQueries.GetSnippetById(r.Context(), id)
	if err != nil {
		utilites.ResponseWithError(w, r, http.StatusNotFound, "")
		return
	}
	if snippet.UserID != user.ID {
		utilites.ResponseWithError(w, r, http.StatusUnauthorized, "")
		return
	}
	s.DbQueries.DeleteSnippetById(r.Context(), snippet.ID)
	utilites.ResponseWithJson(w, r, http.StatusNoContent, "")

}
