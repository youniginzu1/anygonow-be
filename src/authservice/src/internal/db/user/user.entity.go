package user

import (
	"github.com/aqaurius6666/authservice/src/internal/db/role"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
)

type User struct {
	database.BaseModel
	Username            *string `gorm:"index:unique"`
	PublicKey           *string `gorm:"index:unique"`
	EncryptedPrivateKey *string
	Mail                *string `gorm:"index:unique"`
	Phone               *string
	RoleID              uuid.UUID
	Role                *role.Role `gorm:"foreignKey:RoleID"`
	IsActive            *bool      `gorm:"default:true"`
	IsDefaultPassword   *bool      `gorm:"type:bool;default:false"`
	RefCode             *string    `gorm:"-:migration;-"`
}

type Search struct {
	database.DefaultSearchModel
	User
}
