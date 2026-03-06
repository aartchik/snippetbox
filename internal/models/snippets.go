package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)



type SnippetModelInterface interface {
    Insert(title string, content string, expires, user_id int) (int, error) 
    Get(snip_id, user_id int) (*Snippet, error)
    Latest(user_id int) ([]*Snippet, error)
	Delete(snippet_id int) (error)
	Update(title string, content string, expires, snippet_id int) (error)
	GetSearch(title string, user_id int) ([]*Snippet, error)
}

type SnippetModelCacheInterface interface {
	GetCache(ctx context.Context, key string) (string, error)
	SetCache(ctx context.Context, key, value string, ttl time.Duration) error
	DelCache(ctx context.Context, keys ...string) error 
}



type Snippet struct {
	ID      int 	  `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
	User_id int		  `json:"user_id"`
}

type SnippetModel struct {
	DB *sql.DB
	RDB *redis.Client
}

func (r *SnippetModel) GetCache(ctx context.Context, key string) (string, error) {
	val, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (r *SnippetModel) SetCache(ctx context.Context, key, value string, ttl time.Duration) error {
	err := r.RDB.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *SnippetModel) DelCache(ctx context.Context, keys ...string) error {
	err := r.RDB.Del(ctx, keys...).Err()
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

	ctx := context.Background()
	key := fmt.Sprintf("snippet:%d", snip_id)

	res, err := m.GetCache(ctx, key)
	if res != "" { 
		var s Snippet
		err = json.Unmarshal([]byte(res), &s)
		if err == nil {
			return &s, nil
		}
	}
	
	stmt := `Select id, title, content, created, expires, user_id from snippets where id = ? and expires > NOW() and user_id = ?`
	row := m.DB.QueryRow(stmt, snip_id, user_id)
	s := &Snippet{}
	err = row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.User_id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	js, err :=  json.Marshal(s)
	if err == nil {
		_ = m.SetCache(ctx, key, string(js), 15 * time.Second)
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
		err = rows.Scan(&s.ID,&s.Title, &s.Content, &s.Created, &s.Expires, &s.User_id)
		
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


func (m *SnippetModel) Delete(snippet_id int) (error) {
	key := fmt.Sprintf("snippet:%d", snippet_id)
	err := m.DelCache(context.Background(), key)
	
	stmt := `delete from snippets where id = ?`
	_, err = m.DB.Exec(stmt, snippet_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *SnippetModel) Update(title string, content string, expires, snippet_id int) (error) {
	key := fmt.Sprintf("snippet:%d", snippet_id)
	err := m.DelCache(context.Background(), key)

	stmt := `update snippets set title = ?, content = ?, expires = DATE_ADD(NOW(), INTERVAL ? DAY)  where id = ?`
	_, err = m.DB.Exec(stmt, title, content, expires, snippet_id)
	if err != nil {
		return err
	}
	return nil

}


func (m *SnippetModel) GetSearch(title string, user_id int) ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires, user_id,
			MATCH(title, content) AGAINST (? IN BOOLEAN MODE) AS score
			FROM snippets
			WHERE expires > NOW() AND user_id = ? AND MATCH(title, content) AGAINST (? IN BOOLEAN MODE)
			ORDER BY score DESC, created DESC
			LIMIT 50;`
	
	q := strings.TrimSpace(title)
	q = q + "*"
	rows, err := m.DB.Query(stmt, q, user_id, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		var score float64
		err = rows.Scan(&s.ID,&s.Title, &s.Content, &s.Created, &s.Expires, &s.User_id, &score)
		
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