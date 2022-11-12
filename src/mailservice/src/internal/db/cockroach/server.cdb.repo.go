package cockroach

import (
	"time"

	"github.com/aqaurius6666/go-utils/database/cockroach"
	"github.com/aqaurius6666/mailservice/src/internal/db/device"
	"github.com/google/wire"
)

// var (
// 	_ db.ServerRepo = (*ServerCDBRepo)(nil)
// )
var (
	timeout = 2 * time.Second
)
var CDBRepoSet = wire.NewSet(wire.Struct(new(ServerCDBRepo), "*"), InterfacesProvider, wire.Struct(new(cockroach.CDBRepository), "*"))

func InterfacesProvider() cockroach.DBInterfaces {
	return cockroach.DBInterfaces{
		device.Device{},
	}
}

type ServerCDBRepo struct {
	cockroach.CDBRepository
}
