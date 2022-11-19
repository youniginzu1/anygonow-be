package state

import "context"

type StateRepo interface {
	SelectState(context.Context, *Search) (*State, error)
	InsertState(context.Context, *State) (*State, error)
	ListStates(context.Context, *Search) ([]*State, error)
}
