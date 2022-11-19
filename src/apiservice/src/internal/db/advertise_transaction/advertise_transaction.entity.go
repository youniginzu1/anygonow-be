package advertise_transaction

import (
	"github.com/aqaurius6666/go-utils/database"
)

type AdvertiseTransaction struct {
	database.BaseModel
	Name         *string
	Price        *float64
	BannerUrl    *string
	Description  *string
	Zipcode      *string
	CategoryName *string
}

type Search struct {
	database.DefaultSearchModel
	AdvertiseTransaction
}
