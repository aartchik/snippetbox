package models

import (
	"database/sql"
	"errors"  // New import
	"strings" // New import
	"time"

	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"     // New import
)

type UserModelInterface interface {
    Insert(name, email, password string) error
    Authenticate(email, password string) (int, error)
    Exist(id int) (bool, error)
	ReturnCorrectPassword(password string, user_id int) (bool, error)
	ChangePassword(password string, user_id int) (error)
	SamePassword(new_password, confirm_password string) (bool)
	ReturnData(id int) (*User, error)
}

type User struct {
	ID int 
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT into users(name, email, hashed_password, created) values (?, ?, ?, NOW())`

	_, err = m.DB.Exec(stmt, name, email, hash)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
                return ErrDuplicateEmail
            }
        }
        return err
    }
	return nil 
}


func (m *UserModel) ChangePassword(password string, user_id int) (error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set hashed_password = ? where id = ?`

	_, err = m.DB.Exec(stmt, hash, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) SamePassword(new_password, confirm_password string) (bool) {
	return new_password == confirm_password
}

func (m *UserModel) ReturnCorrectPassword(password string, user_id int) (bool, error) {
	var hash []byte
	stmt := `select hashed_password from users where id = ?`

	err := m.DB.QueryRow(stmt, user_id).Scan(&hash)
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

func (m *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hash []byte
	stmt := `select id, hashed_password from users where email = ?`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hash)
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

func (m *UserModel) Exist(id int) (bool, error) {
	var exists bool

	stmt := "select exists(select true from users where id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
    return exists, err
}

func (m *UserModel) ReturnData(id int) (*User, error) {
	stmt := "select name, email, created from users where id = ?"
	row := m.DB.QueryRow(stmt, id)
	s := &User{}
	err := row.Scan(&s.Name, &s.Email, &s.Created)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}


