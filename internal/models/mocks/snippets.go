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

func (m *SnippetModel) Insert(title string, content string, expires, user_id int) (int, error) {
    return 2, nil
}

func (m *SnippetModel) Get(snip_id, user_id int) (*models.Snippet, error) {
    switch snip_id {
    case 1:
        return mockSnippet, nil
    default:
        return nil, models.ErrNoRecord
    }
}

func (m *SnippetModel) Latest(user_id int) ([]*models.Snippet, error) {
    return []*models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) Delete(snippet_id int) (error) {
    return nil
}
func (m *SnippetModel) Update(title string, content string, expires, snippet_id int) (error) {
    return nil
}