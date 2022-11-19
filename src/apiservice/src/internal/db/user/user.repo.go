package user

import "context"

type UserRepo interface {
	SelectUser(context.Context, *Search) (*User, error)
	InsertUser(context.Context, *User) (*User, error)
	ListUsers(context.Context, *Search) ([]*User, error)
	UpdateUser(context.Context, *Search, *User) error
	TotalUsers(context.Context, *Search) (*int64, error)
	DeleteUser(context.Context, *Search) error
	SelectUserProcess(context.Context, *Search) (*User, error)
}
