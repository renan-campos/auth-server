/*
This example client will get a JWT from the authentication server, print
out its contents, then use the JWT to interact with a different example
service.
*/
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	jwtConsts "github.com/renan-campos/auth-server/pkg/jwt"
	"github.com/renan-campos/auth-server/pkg/otp"
)

func main() {
	otpFn := flag.String("otp-secret-file", "", "The name of the file holding the otp secret")
	flag.Parse()
	if otpFn == nil || *otpFn == "" {
		panic("The --otp-secret-file must be passed")
	}
	authenticator, err := otp.NewAuthenticator(*otpFn)
	if err != nil {
		log.Fatal(err)
	}

	// Make a POST request
	fmt.Println("\nMaking POST request...")
	resp, err := http.Post(
		"http://localhost:8008/token",
		"application/json",
		bytes.NewBuffer([]byte(authenticator.CurrentToken())))
	if err != nil {
		log.Fatal("POST request failed:", err)
	}
	defer resp.Body.Close()

	// Read response body
	rawJwt, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}
	fmt.Printf("Raw JSON Web Token:\n===\n%s\n===\n", rawJwt)

	jwtParts := strings.Split(string(rawJwt), ".")
	if len(jwtParts) != 3 {
		log.Fatal("Expected there to be three parts to the JSON Web token")
	}
	decode := func(str string) string {
		decodedStr, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			log.Fatal("Decoder failure:", err)
		}
		return string(decodedStr)
	}

	fmt.Printf("Manually Decoded JSON Web Token Header:\n===\n%s\n===\n", decode(jwtParts[0]))
	// I cannot explain the additional == characters, somehow I just knew to do it...
	fmt.Printf("Manually Decoded JSON Web Token Payload:\n===\n%s\n===\n", decode(jwtParts[1]+"=="))
	fmt.Printf("Encrypted JSON Web Token Signature:\n===\n%s\n===\n", jwtParts[2])

	webToken, err := jwt.ParseSigned(string(rawJwt))
	if err != nil {
		log.Fatal("Failed to parse JWT:", err)
	}

	prettyPrintHeaders := func(headers []jose.Header) string {
		var outStr string
		for headerNum, header := range headers {
			outStr += fmt.Sprintf("[Header #%d start]\n\tAlgorithm: %s\n\tType: %s\n[Header #%d end]",
				headerNum+1,
				header.Algorithm,
				header.ExtraHeaders[jose.HeaderType],
				headerNum+1,
			)
		}
		return outStr
	}

	prettyPrintClaims := func(claims jwt.Claims) string {
		// Load the Eastern Standard Time (EST) location
		estLocation, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.Fatal("Error loading EST location:", err)
		}

		return fmt.Sprintf("\tAudience: %70s\n\tSubject:  %70s\n\tIssuer:   %70s\n\tIssued_At: %70s\n\tExpiry:    %70s",
			strings.Join(claims.Audience, " , "),
			claims.Subject,
			claims.Issuer,
			claims.IssuedAt.Time().In(estLocation).Format("Monday, January 2, 2006 03:04:05 PM EST"),
			claims.Expiry.Time().In(estLocation).Format("Monday, January 2, 2006 03:04:05 PM EST"),
		)
	}

	fmt.Printf("JSON Web Token Header:\n===\n%s\n===\n", prettyPrintHeaders(webToken.Headers))
	var unverifiedClaims jwt.Claims
	err = webToken.UnsafeClaimsWithoutVerification(&unverifiedClaims)
	if err != nil {
		log.Fatal("Failed to get JWT claims:", err)
	}
	fmt.Printf("JSON Web Token Payload (Unverified):\n===\n%v\n===\n", prettyPrintClaims(unverifiedClaims))
	err = unverifiedClaims.Validate(jwt.Expected{
		Issuer: jwtConsts.Issuer,
		Time:   time.Now(),
	})
	if err != nil {
		log.Fatal("Failed to validate claims:", err)
	}
	log.Println("Claims validated sucessful.\n")

	var verifiedClaims jwt.Claims
	resp, err = http.Get("http://localhost:8008/jwks")
	if err != nil {
		log.Fatal("Http request failed:", err)
	}
	defer resp.Body.Close()

	var jwks jose.JSONWebKeySet
	marshalledJwks, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read GET response body")
	}
	if err := json.Unmarshal(marshalledJwks, &jwks); err != nil {
		log.Fatal("Failed to unmarshall jwks")
	}

	err = webToken.Claims(jwks.Keys[0], &verifiedClaims)
	if err != nil {
		log.Fatal("Failed to verify jwt: ", err)
	}
	log.Println("JWT verified!")

	// Sending an http message with a jwt
	req, err := http.NewRequest(http.MethodGet,
		"http://localhost:8008/token",
		bytes.NewBuffer([]byte(authenticator.CurrentToken())))
	if err != nil {
		log.Fatal("Http request failed:", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", rawJwt))
	fmt.Println(req)
}
