package auth

import (
	"github.com/saas-be-usergroup/internal/core/domain/mailer"
	"github.com/xlzd/gotp"
	"strconv"
	"strings"
	"time"
)

type OTP struct {
	UUID     string
	OTP      string
	Duration uint64
}

func NewOTP(uuid string, otp string) *OTP {
	return &OTP{
		UUID: uuid,
		OTP:  otp,
	}
}

func (o OTP) GetUUID() string {
	return o.UUID
}

func (o OTP) GetOTP() string {
	return o.OTP
}

func (o *OTP) SetDuration(duration uint64) {
	o.Duration = duration
}

func (o *OTP) IsNotFound() bool {
	return o == nil || o.OTP == ""
}

func GenerateNewOTP() string {
	return gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO").Now()
}

func SetExpiredInSecond(expiredInSecond uint64) (*time.Duration, error) {
	inSecond, err := time.ParseDuration(strconv.FormatUint(expiredInSecond, 10) + "s")
	if err != nil {
		return nil, err
	}

	return &inSecond, nil
}

func (o OTP) ToRegisterMail() *mailer.RegisterMail {
	return &mailer.RegisterMail{
		OTP:      o.OTP,
		Duration: o.Duration,
	}
}

func (o OTP) ToRegisterResponse() *RegisterBeforeResponse {
	return &RegisterBeforeResponse{UUID: o.UUID}
}

func (o OTP) IsValid(otpRequest string) bool {
	return strings.EqualFold(o.OTP, otpRequest)
}
