package lib

import (
	"github.com/aqaurius6666/mailservice/src/internal/var/e"
	"github.com/google/uuid"
)

func ToUUID(id interface{}) (uuid.UUID, error) {
	switch tmp := id.(type) {
	case uuid.UUID:
		return tmp, nil
	case string:
		tmp2, err := uuid.Parse(tmp)
		if err != nil {
			return uuid.Nil, e.ErrIdInvalidFormat
		}
		return tmp2, nil
	default:
		return uuid.Nil, e.ErrIdInvalidFormat
	}
}
