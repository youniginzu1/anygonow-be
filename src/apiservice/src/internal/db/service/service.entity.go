package service

import (
	"github.com/aqaurius6666/apiservice/src/internal/db/category"
	"github.com/aqaurius6666/go-utils/database"
	"github.com/google/uuid"
	"github.com/ubgo/gormuuid"
)

type Service struct {
	database.BaseModel
	Name        *string            `gorm:"-:migration;->"`
	NumberOrder *int64             `gorm:"-:migration;->"`
	LogoUrl     *string            `gorm:"-:migration;->"`
	Status      *int32             `gorm:"type:int8;default:0"`
	BusinessId  uuid.UUID          `gorm:"type:uuid"`
	CategoryId  uuid.UUID          `gorm:"type:uuid"`
	Category    *category.Category `gorm:"foreignKey:CategoryId"`
}

type Search struct {
	database.DefaultSearchModel
	Service
	CategoryIds gormuuid.UUIDArray
}
