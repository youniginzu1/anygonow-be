package twilloclient

type TwilloMessageCallbackData struct {
	SmsMessageSid    string
	NumMedia         int
	ToCity           string
	FromZip          string
	SmsSid           string
	FromState        string
	SmsStatus        string
	FromCity         string
	Body             string
	FromCountry      string
	To               string
	ToZip            string
	NumSegements     int
	ReferralNumMedia int
	MessageSid       string
	AccountSid       string
	From             string
	ApiVersion       string
	ToCountry        string
}
