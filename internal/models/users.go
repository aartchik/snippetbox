package models

import (
	"database/sql"
	"errors"  // New import
	"strings" // New import
	"time"

	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"     // New import
)

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


