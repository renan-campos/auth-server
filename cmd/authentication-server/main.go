/*
The authentication server verifies a client's identity and returns a JSON
Web Token (JWT).
(In this case, the identity will simply always be 'guest')

The JWT will be passed from the client to the service that needs to know
the identity of the client. The service will query this authentication
server for the JSON Web Key Set (JWKS) to verify the JWT that is passed.

There are conventions for both the JWT issuance endpoint and the JWKS
endpoint, although they may vary based on the implementation and specific
requirements of the authentication server.
*/
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	otpFn := flag.String("otp-secret-file", "", "The name of the file holding the otp secret")
	flag.Parse()

	if otpFn == nil || *otpFn == "" {
		panic("The --otp-secret-file must be passed")
	}
	authenticator, err := NewOtpAuthenticator(*otpFn)
	if err != nil {
		log.Fatal(err)
	}

	jsonWebKey := generateJsonWebKey()

	// Define the HTTP handlers
	http.HandleFunc(
		"/token",
		generateTokenIssuer(
			webKeyToSigner(jsonWebKey),
			authenticator,
		),
	)

	http.HandleFunc(
		"/jwks",
		generateJsonWebKeysRequestHandler(
			createKeySet(jsonWebKey),
		),
	)

	// Start the HTTP server
	log.Println("Athentication Server started on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
