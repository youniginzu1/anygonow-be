package c

import "time"

var (
	OTP_EXPIRE_TIME = 5 * time.Minute
)

const (
	OTP_TYPE_REGISTER = iota
)

const (
	OTP_STATUS_PENDING = iota
)

const (
	ROLE_CUSTOMER = iota
	ROLE_HANDYMAN
	ROLE_ADMIN
)

var (
	SERVICE_NAME       = "mail-service"
	ID           int64 = 5
)
