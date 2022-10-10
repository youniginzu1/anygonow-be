package c

import "time"

var (
	OTP_EXPIRE_TIME = 5 * time.Minute
	// OTP_SPAM_TIME         = 1 * time.Minute
	OTP_SPAM_TIME       = 10 * time.Second
	SERVICE_NAME        = "auth-service"
	ID            int64 = 1
	OTP_REGISTER_TIME = 365 * 24 * time.Hour
)
