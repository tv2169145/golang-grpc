package repos

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/tv2169145/golang-grpc/types"
)

// UsersRepo - users repo interface
type UsersRepo interface {
	Create(*types.User) error
	FindById(int64) (*types.User, error)
	FindByEmail(string) (*types.User, error)
}

// NewUsersRepo - return a new users repo
func NewUsersRepo(db *xorm.Engine) UsersRepo {
	return &usersRepo{db:db}
}

type usersRepo struct {
	db *xorm.Engine
}

func (u *usersRepo) Create(user *types.User) (err error) {
	if err = types.Validate(user); err != nil {
		return
	}
	if _, err = u.db.Insert(user); err != nil {
		return
	}
	return nil
}

func (u *usersRepo) FindById(id int64) (user *types.User, err error) {
	if id <= 0 {
		err = fmt.Errorf("invalid user id")
		return
	}
	user = new(types.User)
	has, err := u.db.ID(id).Get(user)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("unable to find a user")
		return
	}
	return
}

func (u *usersRepo) FindByEmail(email string) (user *types.User, err error) {
	if len(email) == 0 {
		err = fmt.Errorf("invalid user email")
		return
	}
	user = new(types.User)
	user.Email = email
	
	has, err := u.db.Get(user)
	if err != nil {
		return 
	}
	if !has {
		err = fmt.Errorf("unable to find a user")
		return 
	}
	return
}
