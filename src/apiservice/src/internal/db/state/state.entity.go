package state

import "github.com/aqaurius6666/go-utils/database"

type State struct {
	database.BaseModel
	Name  *string
	Short *string
}

type Search struct {
	database.DefaultSearchModel
	State
}
