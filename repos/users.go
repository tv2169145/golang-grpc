package repos

import "github.com/go-xorm/xorm"

// UsersRepo - users repo interface
type UsersRepo interface {

}

type usersRepo struct {
	db *xorm.Engine
}

// NewUsersRepo - return a new users repo
func NewUsersRepo(db *xorm.Engine) UsersRepo {
	return &usersRepo{db:db}
}
