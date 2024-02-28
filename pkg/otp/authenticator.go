package otp

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/xlzd/gotp"
)

func NewAuthenticator(secretFileName string) (*Authenticator, error) {
	fp, err := os.Open(secretFileName)
	if err != nil {
		return nil, err
	}

	out, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	// Note: Only the first 32 bytes of the file matters
	return &Authenticator{
		totp: gotp.NewDefaultTOTP(string(out[:32])),
	}, nil
}

func (a *Authenticator) VerifyHttpRequest(r *http.Request) (valid bool, err error) {
	otp, err := io.ReadAll(r.Body)
	if err != nil {
		return valid, err
	}
	return a.totp.VerifyTime(string(otp), time.Now()), nil
}

func (a *Authenticator) CurrentToken() string {
	return a.totp.Now()
}
