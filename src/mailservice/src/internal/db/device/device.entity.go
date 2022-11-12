package device

import (
	"github.com/aqaurius6666/go-utils/database"
)

var (
	ACTIVE   = 0
	INACTIVE = 1
)

type Device struct {
	database.BaseModel
	DeviceId *string
	UserId   *string
	Active   *int `gorm:"default:0;type:int8"` // 0 Active, 1 Inactive
}

type Search struct {
	Device
	database.DefaultSearchModel
}
