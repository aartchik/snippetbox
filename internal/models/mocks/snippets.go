package mocks

import (
	"time"

	"snippetbox.net/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires, userID, visibilityLevel int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(snippetID, userID int) (*models.Snippet, error) {
	switch snippetID {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest(userID int) ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) Delete(snippetID, userID int) error {
	return nil
}

func (m *SnippetModel) Update(title string, content string, expires, snippetID, userID int) error {
	return nil
}

func (m *SnippetModel) GetSearch(title string, userID int) ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
