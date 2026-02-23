package models

import (
	"database/sql"
	"time"
	"errors" 
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
	User_id int
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Delete(snippet_id int) (error) {
	stmt := `delete from snippets where id = ?`
	_, err := m.DB.Exec(stmt, snippet_id)
	if err != nil {
		return err
	}
	return nil

}

func (m *SnippetModel) Insert(title string, content string, expires, user_id int) (int, error) {
	stmt := `INSERT INTO snippets(title, content, created, expires, user_id) VALUES(?, ?, NOW(),  DATE_ADD(NOW(), INTERVAL ? DAY), ?)`

	result, err := m.DB.Exec(stmt, title, content, expires, user_id)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(snip_id, user_id int) (*Snippet, error) {
	stmt := `Select * from snippets where id = ? and expires > NOW() and user_id = ?`

	row := m.DB.QueryRow(stmt, snip_id, user_id)
	s := &Snippet{}
	err := row.Scan(&s.ID, &s.User_id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest(user_id int) ([]*Snippet, error) {
	stmt := `SELECT * from snippets where expires > NOW() and user_id = ? order by created desc limit 10`

	rows, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID,&s.User_id, &s.Title, &s.Content, &s.Created, &s.Expires)
		
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	  if err = rows.Err(); err != nil {
        return nil, err
    }
	return snippets, nil
}