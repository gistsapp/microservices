package core

import (
	"database/sql"

	"github.com/gistsapp/api/auth/repositories"
	"github.com/gistsapp/api/types"
)

type UserService interface {
	GetUserByID(id string) (*types.User, error)
	CreateUser(user *types.User) (*types.User, error)
	DeleteUser(id string) error
	UpdateUser(user *types.User) (*types.User, error)
}

type userService struct {
	db repositories.Database
}

func NewUserService(db repositories.Database) UserService {
	return &userService{
		db: db,
	}
}

func (u *userService) GetUserByID(id string) (*types.User, error) {
	user, err := u.db.GetUserByID(id)
	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}else if err != nil {
		return nil, err
	}
	return user, err
}

func (u *userService) CreateUser(user *types.User) (*types.User, error) {
	return u.db.CreateUser(user)
}

func (u *userService) DeleteUser(id string) error {
	err := u.db.DeleteUser(id)
	if err == sql.ErrNoRows {
		return types.ErrNotFound
	}
	return err
}

func (u *userService) UpdateUser(user *types.User) (*types.User, error) {
	user, err := u.db.UpdateUser(user)
	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	return user, err
}

