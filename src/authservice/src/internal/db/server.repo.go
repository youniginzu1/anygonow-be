package db

import (
	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/internal/db/role"
	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/go-utils/database"
)

type DBDsn string

type ServerRepo interface {
	database.CommonRepository
	user.UserRepo
	role.RoleRepo
	otp.OtpRepo
}
