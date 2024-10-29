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
package jwt

import (
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Authenticator interface {
	authenticate(r *http.Request) (valid bool, err error)
}

type HttpVerifier interface {
	VerifyHttpRequest(r *http.Request) (valid bool, err error)
}

func GenerateTokenIssuer(
	signer jose.Signer,
	authenticator HttpVerifier,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		valid, err := authenticator.VerifyHttpRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Error(w, "Invalid OTP", http.StatusUnauthorized)
			return
		}

		cl := jwt.Claims{
			Subject:  Subject,
			Issuer:   Issuer,
			Audience: jwt.Audience{IntendedAudience},
			// Have the token expire in an hour
			Expiry: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			// A second is subtracted from the issed at time, to avoid issued in future errors
			IssuedAt: jwt.NewNumericDate(time.Now()),
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
