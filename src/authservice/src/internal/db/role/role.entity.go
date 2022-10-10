package role

import "github.com/aqaurius6666/go-utils/database"

type Role struct {
	database.BaseModel
	Name *string
	Code *int
}

type Search struct {
	database.DefaultSearchModel
	Role
}
