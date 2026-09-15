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
	Insert(title string, content string, expires, user_id, visibility_level int) (int, error)
	Get(snip_id, user_id int) (*Snippet, error)
	Latest(user_id int) ([]*Snippet, error)
	Delete(snippet_id, user_id int) error
	Update(title string, content string, expires, snippet_id, user_id int) error
	GetSearch(title string, user_id int) ([]*Snippet, error)
}

type SnippetModelCacheInterface interface {
	GetCache(ctx context.Context, key string) (string, error)
	SetCache(ctx context.Context, key, value string, ttl time.Duration) error
	DelCache(ctx context.Context, keys ...string) error
}

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
	User_id int       `json:"user_id"`
}

type SnippetModelCache struct {
	RDB *redis.Client
}

type SnippetModelWithPsql struct {
	DB *sql.DB
	SnippetModelCache
}

func (r *SnippetModelCache) GetCache(ctx context.Context, key string) (string, error) {
	val, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (r *SnippetModelCache) SetCache(ctx context.Context, key, value string, ttl time.Duration) error {
	err := r.RDB.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *SnippetModelCache) DelCache(ctx context.Context, keys ...string) error {
	err := r.RDB.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *SnippetModelWithPsql) Insert(
	title string,
	content string,
	expires, userID, visibilityLevel int,
) (int, error) {

	stmt := `
        INSERT INTO snippets (
            title,
            content,
            created,
            expires,
            user_id,
            visibility_level
        )
        VALUES (
            $1,
            $2,
            NOW(),
            NOW() + ($3 * INTERVAL '1 day'),
            $4,
            $5
        )
        RETURNING id
    `

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	var id int

	err := m.DB.QueryRowContext(
		ctx,
		stmt,
		title,
		content,
		expires,
		userID,
		visibilityLevel,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModelWithPsql) Get(snip_id, user_id int) (*Snippet, error) {

	ctx := context.Background()
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	key := fmt.Sprintf("snippet:%d", snip_id)

	res, err := m.GetCache(ctx, key)
	if res != "" {
		var s Snippet
		err = json.Unmarshal([]byte(res), &s)
		if err == nil {
			return &s, nil
		}
	}

	stmt := `Select id, title, content, created, expires, user_id from snippets where id = $1 and expires > NOW() and user_id = $2`
	row := m.DB.QueryRowContext(ctx2, stmt, snip_id, user_id)
	s := &Snippet{}
	err = row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.User_id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	js, err := json.Marshal(s)
	if err == nil {
		_ = m.SetCache(ctx, key, string(js), 15*time.Second)
	}

	return s, nil
}

func (m *SnippetModelWithPsql) Latest(user_id int) ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires, user_id from snippets where expires > NOW() and user_id = $1 order by created desc, id desc limit 10`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, stmt, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires, &s.User_id)

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

func (m *SnippetModelWithPsql) Delete(snippet_id, user_id int) error {
	stmt := `delete from snippets where id = $1 and user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, stmt, snippet_id, user_id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRecord
	}

	key := fmt.Sprintf("snippet:%d", snippet_id)
	err = m.DelCache(context.Background(), key)
	if err != nil {
		return err
	}

	return nil
}

func (m *SnippetModelWithPsql) Update(title string, content string, expires, snippet_id, user_id int) error {

	stmt := `update snippets set title = $1, content = $2, expires =  NOW() + ($3 * INTERVAL '1 day') where id = $4 and user_id = $5`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := m.DB.ExecContext(ctx, stmt, title, content, expires, snippet_id, user_id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrNoRecord
	}

	key := fmt.Sprintf("snippet:%d", snippet_id)
	m.DelCache(context.Background(), key)
	return nil

}

func (m *SnippetModelWithPsql) GetSearch(
	title string,
	userID int,
) ([]*Snippet, error) {

	stmt := `
        SELECT
            id,
            title,
            content,
            created,
            expires,
            user_id,
            ts_rank(
                to_tsvector('english', title || ' ' || content),
                plainto_tsquery('english', $1)
            ) AS score
        FROM snippets
        WHERE
            expires > NOW()
            AND user_id = $2
            AND to_tsvector('english', title || ' ' || content)
                @@ plainto_tsquery('english', $1)
        ORDER BY score DESC, created DESC
        LIMIT 50
    `

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	q := strings.TrimSpace(title)

	rows, err := m.DB.QueryContext(ctx, stmt, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		var score float64

		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Content,
			&s.Created,
			&s.Expires,
			&s.User_id,
			&score,
		)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
