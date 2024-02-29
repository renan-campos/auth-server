package main

import (
	"fmt"
	"log"
	"time"

	"github.com/renan-campos/auth-server/pkg/jwt"
)

func main() {
	jwk := jwt.GenerateJsonWebKey()
	jwks := jwt.ExtractPublicJsonWebKey(jwk)

	signer := jwt.WebKeyToSigner(jwk)
	claims := jwt.NewClaims(jwt.ClaimData{
		Subject:  "test-subject",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		// Have the token expire in an hour
		Expiry: time.Now().Add(time.Hour),
		// A second is subtracted from the issed at time, to avoid issued in future errors
		IssuedAt: time.Now(),
	})
	token, err := jwt.CreateToken(signer, claims)
	if err != nil {
		log.Fatal("Failed to sign a token: ", err)
	}

	jsonJwks, err := jwks.MarshalJSON()
	if err != nil {
		log.Fatal("Failed to marshal jwks into json: ", err)
	}

	fmt.Printf(`{
        "jwks": %q,
        "jwt": %q
}`, string(jsonJwks), token)

}
