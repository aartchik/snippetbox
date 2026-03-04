package models

import (
    "testing"

    "snippetbox.net/internal/assert"
    "time"
)

func TestSnippetModelExist(t *testing.T) {
    tests := []struct {
        name   string
        userID int
        want   bool
    }{
        {
            name:   "Valid ID",
            userID: 1,
            want:   true,
        },
        {
            name:   "Zero ID",
            userID: 0,
            want:   false,
        },
        {
            name:   "Non-existent ID",
            userID: 2,
            want:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db := newTestDB(t)

            m := UserModel{db}

            exists, err := m.Exist(tt.userID)

            assert.Equal(t, exists, tt.want)
            assert.NilError(t, err)
        })
    }
}


func TestSnippetModelUpdate(t *testing.T) {
	db := newTestDB(t)
	m := SnippetModel{DB: db}

	id, err := m.Insert("old title", "old content", 7, 1)
	assert.NilError(t, err)

    start := time.Now().UTC()
	err = m.Update("new title", "new content", 1, id)
	assert.NilError(t, err)

	s, err := m.Get(id, 1)
	assert.NilError(t, err)

	assert.Equal(t, s.Title, "new title")
	assert.Equal(t, s.Content, "new content")

    if !(s.Expires.After(start) || s.Expires.Before(start.Add(24*time.Hour + 5 * time.Minute))) {
        t.Errorf("got: %v; beetween: %v and %v", s.Expires, start.UTC(), start.Add(24*time.Hour + 5 * time.Minute))
    }
}

func TestSnippetModelDelete(t *testing.T) {
    db := newTestDB(t)
    m := SnippetModel{DB: db}
    var cnt, cnt_new int

    result := m.DB.QueryRow("select count(*) from snippets")
    err := result.Scan(&cnt)
    assert.NilError(t, err)
    
    id, err := m.Insert("old title", "old content", 7, 1)
	assert.NilError(t, err)

    err = m.Delete(id)
    assert.NilError(t, err)

    result = m.DB.QueryRow("select count(*) from snippets")
    err = result.Scan(&cnt_new)

    assert.Equal(t, cnt_new, cnt)

}