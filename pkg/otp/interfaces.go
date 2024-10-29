package otp

import (
	"net/http"
)

type Authenticator interface {
	// VerifyHttpRequest takes as input a request body expected to have the
	// one-time-password in its request body. Returned is a boolean to
	// indicate if the otp is valid at the current time.
	VerifyHttpRequest(r *http.Request) (valid bool, err error)

	// CurrentToken
	CurrentToken() string
}
