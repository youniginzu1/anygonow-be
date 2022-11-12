package db

import (
	"github.com/aqaurius6666/go-utils/database"
	"github.com/aqaurius6666/mailservice/src/internal/db/device"
)

type DBDsn string

type ServerRepo interface {
	database.CommonRepository
	device.DeviceRepo
}
