package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/renan-campos/auth-server/pkg/jwt"
)

func main() {
	keyJsonFn := flag.String("private-key", "/tmp/token-generator-key",
		"The location of the JSON private key. Set to empty for private key initiation.")
	fIssuer := flag.String("issuer", "", "The issuer url of the token")
	fAudience := flag.String("audience", "", "The audience of the token")
	fSubject := flag.String("subject", "", "The subject of the token")

	flag.Parse()

	if keyJsonFn == nil || *keyJsonFn == "" {
		log.Print("Key creation initiated")
		jwk := jwt.GenerateJsonWebKey()
		data, err := jwk.MarshalJSON()
		if err != nil {
			log.Fatal("Failed to marshal key into json: ", err)
		}
		err = os.WriteFile("/tmp/token-generator-key", data, 0666)
		if err != nil {
			log.Fatal("Failed to create generator-key: ", err)
		}
		log.Print("Key creation complete")

		*keyJsonFn = "/tmp/token-generator-key"
	}

	checkRequiredStringArg := func(arg *string, argName string) {
		if arg == nil || *arg == "" {
			log.Fatalf("command line variable %q required", argName)
		}
	}
	checkRequiredStringArg(fIssuer, "issuer")
	checkRequiredStringArg(fAudience, "audience")
	checkRequiredStringArg(fSubject, "subject")

	jwk, err := jwt.KeyFromJson(*keyJsonFn)
	if err != nil {
		log.Fatal("Failed to read generator key: ", err)
	}

	jwks := jwt.CreateKeySet(jwt.ExtractPublicJsonWebKey(jwk))

	signer := jwt.WebKeyToSigner(jwk)
	claims := jwt.NewClaims(jwt.ClaimData{
		Issuer:   *fIssuer,
		Subject:  *fSubject,
		Audience: *fAudience,
		// Have the token expire in a week
		Expiry: time.Now().Add(1 * time.Hour),
		// A second is subtracted from the issed at time, to avoid issued in future errors
		IssuedAt: time.Now(),
	})
	token, err := jwt.CreateToken(signer, claims)
	if err != nil {
		log.Fatal("Failed to sign a token: ", err)
	}

	jsonJwks, err := json.Marshal(jwks)
	if err != nil {
		log.Fatal("Failed to marshal jwks into json: ", err)
	}

	fmt.Printf(`{
        "jwks": %q,
        "jwt": %q
}`, string(jsonJwks), token)

}
