package c

import "time"

// const (
// 	OTP_TYPE_REGISTER = iota
// 	OTP_TYPE_FORGOT_PASSWORD
// )

// const (
// 	SERVICE_STATUS_ACTIVE = iota
// 	SERVICE_STATUS_INACTIVE
// )

const (
	PRESIGNED_URL_EXPIRE_TIME = 5 * time.Minute
)

var (
	SERVICE_NAME              = "api-service"
	ID                  int64 = 2
	BUY_ADVERTISE_LIMIT int64 = 100
)

// const (
// 	ROLE_CUSTOMER = iota
// 	ROLE_HANDYMAN
// 	ROLE_ADMIN
// )

var (
	REQUEST_HANDYMAN_NOTIFICATION  = "request-notification"
	CANCEL_HANDYMAN_NOTIFICATION   = "cancel-notification"
	COMPLETE_CUSTOMER_NOTIFICATION = "complete-notification"
	FEE_HANDYMAN_NOTIFICATION      = "fee-notification"

	CONNECT_CUSTOMER_NOTIFICATION = "connect-notification"
	REJECT_CUSTOMER_NOTIFICATION  = "reject-notification"
)
