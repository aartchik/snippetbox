package mocks

import (
    "snippetbox.net/internal/models"
)


type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
    switch email {
    case "1@1":
        return models.ErrDuplicateEmail
    default:
        return nil
    }
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
    if email == "1@1" && password == "1234" {
        return 1, nil
    }

    return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exist(id int) (bool, error) {
    switch id {
    case 1:
        return true, nil
    default:
        return false, nil
    }
}


func (m *UserModel) ReturnData(id int) (*models.User, error) {


	return &models.User{}, nil
}


func (m *UserModel) ChangePassword(password string, user_id int) (error) {
	return nil
}

func (m *UserModel) SamePassword(new_password, confirm_password string) (bool) {
	return new_password == confirm_password
}

func (m *UserModel) ReturnCorrectPassword(password string, user_id int) (bool, error) {
	return true, nil
}