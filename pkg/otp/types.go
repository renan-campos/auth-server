package otp

import (
	"net/http"

	"github.com/xlzd/gotp"
)

type Authenticator struct {
	totp *gotp.TOTP
}

type HttpVerifier interface {
	VerifyHttpRequest(r *http.Request) (valid bool, err error)
}
