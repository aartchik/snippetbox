package models

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exist(id int) (bool, error)
	ReturnCorrectPassword(password string, user_id int) (bool, error)
	ChangePassword(password string, user_id int) error
	ReturnData(id int) (*User, error)
	InsertPhoto(id int, avatar_URL string) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	AvatarURL      string
}

type UserModelWithPsql struct {
	DB *sql.DB
}

func randomAvatar() (string, error) {
	dir := "ui/static/img/avatars"
	files, err := os.ReadDir(dir)
	if err != nil {
		return "/static/img/avatars/penguin.png", err
	}

	var images []string

	for _, f := range files {
		fullPath := filepath.Join("/static/img/avatars", f.Name())
		images = append(images, fullPath)
	}
	rand.Seed(time.Now().UnixNano())
	return images[rand.Intn(len(images))], nil
}

func (m *UserModelWithPsql) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	defaultPhoto, err := randomAvatar()
	if err != nil {
		return err
	}
	stmt := `INSERT into users(name, email, hashed_password, created, avatar_url) values ($1, $2, $3, NOW(), $4)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, stmt, name, email, hash, defaultPhoto)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code == "23505" && pgError.Constraint == "users_uc_email" {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (m *UserModelWithPsql) InsertPhoto(id int, avatar_URL string) error {
	stmt := "UPDATE users SET avatar_url=$1 WHERE id=$2"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, stmt, avatar_URL, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModelWithPsql) ChangePassword(password string, user_id int) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set hashed_password = $1 where id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, stmt, hash, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModelWithPsql) ReturnCorrectPassword(password string, user_id int) (bool, error) {
	var hash []byte
	stmt := `select hashed_password from users where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, stmt, user_id).Scan(&hash)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, ErrInvalidCredentials
		} else {
			return false, err
		}
	}
	return true, nil
}

func (m *UserModelWithPsql) Authenticate(email, password string) (int, error) {

	var id int
	var hash []byte
	stmt := `select id, hashed_password from users where email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, stmt, email).Scan(&id, &hash)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil

}

func (m *UserModelWithPsql) Exist(id int) (bool, error) {
	var exists bool
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "select exists(select true from users where id = $1)"
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModelWithPsql) ReturnData(id int) (*User, error) {
	stmt := "select name, email, created, avatar_url from users where id = $1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, stmt, id)
	s := &User{}
	err := row.Scan(&s.Name, &s.Email, &s.Created, &s.AvatarURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
