/*
JWT Issuance Endpoint:

  - This endpoint is typically responsible for generating and issuing JWTs to
    clients after successful authentication.

  - Conventionally, this endpoint is often named /token, /oauth/token, or
    /authenticate.

  - It usually accepts HTTP POST requests and expects credentials or authorization
    codes in the request body or as URL parameters, depending on the authentication
    flow being used (e.g., password grant, authorization code grant, etc.). Upon
    successful authentication and authorization, this endpoint responds with a JWT
    in the response body, along with additional information such as token type,
    expiration time, etc.
*/
package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/xlzd/gotp"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Authenticator interface {
	authenticate(r *http.Request) (valid bool, err error)
}

func generateTokenIssuer(
	signer jose.Signer,
	authenticator Authenticator,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		valid, err := authenticator.authenticate(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Error(w, "Invalid OTP", http.StatusUnauthorized)
			return
		}

		cl := jwt.Claims{
			Subject:  "demo-token",
			Issuer:   "authentication-server",
			Audience: jwt.Audience{"demo-client"},
			// Have the token expire in an hour
			Expiry: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		}
		token, err := jwt.Signed(signer).Claims(cl).CompactSerialize()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		// Send a success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
	}
}

type otpAuthenticator struct {
	totp *gotp.TOTP
}

func NewOtpAuthenticator(secretFileName string) (*otpAuthenticator, error) {
	fp, err := os.Open(secretFileName)
	if err != nil {
		return nil, err
	}

	out, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	// Note: Only the first 32 bytes of the file matters
	return &otpAuthenticator{
		totp: gotp.NewDefaultTOTP(string(out[:32])),
	}, nil
}

func (a *otpAuthenticator) authenticate(r *http.Request) (valid bool, err error) {
	otp, err := io.ReadAll(r.Body)
	if err != nil {
		return valid, err
	}
	return a.totp.VerifyTime(string(otp), time.Now()), nil
}
