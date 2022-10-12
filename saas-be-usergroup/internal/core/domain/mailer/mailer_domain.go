package mailer

type Mailer struct {
	Template  string
	Recipient string
	Subject   string
	Prop      interface{}
}

type RegisterMail struct {
	OTP      string
	Duration uint64
}

func NewRegisterOTPMailer(recipient string, prop interface{}) *Mailer {
	return &Mailer{
		Template:  "register-otp.html",
		Recipient: recipient,
		Subject:   "Register OTP",
		Prop:      prop,
	}
}
