package cockroach

import (
	"github.com/aqaurius6666/authservice/src/internal/db/otp"
	"github.com/aqaurius6666/authservice/src/internal/db/role"
	"github.com/aqaurius6666/authservice/src/internal/db/user"
	"github.com/aqaurius6666/go-utils/database/cockroach"
	"github.com/google/wire"
)

// var (
// 	_ db.ServerRepo = (*ServerCDBRepo)(nil)
// )
var CDBRepoSet = wire.NewSet(wire.Struct(new(ServerCDBRepo), "*"), InterfacesProvider, wire.Struct(new(cockroach.CDBRepository), "*"))

func InterfacesProvider() cockroach.DBInterfaces {
	return cockroach.DBInterfaces{
		user.User{},
		role.Role{},
		otp.Otp{},
	}
}

type ServerCDBRepo struct {
	cockroach.CDBRepository
}
