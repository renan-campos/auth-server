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

	"github.com/renan-campos/auth-server/pkg/jwt"
	"github.com/renan-campos/auth-server/pkg/otp"
)

func main() {

	otpFn := flag.String("otp-secret-file", "", "The name of the file holding the otp secret")
	assetsDir := flag.String("assets-dir", "", "The name of the directory containing html assets")
	flag.Parse()

	// Check flags {
	if otpFn == nil || *otpFn == "" {
		panic("The --otp-secret-file must be passed")
	}
	if assetsDir == nil || *assetsDir == "" {
		panic("The assets directory must be specified")
	}
	// }

	authenticator, err := otp.NewAuthenticator(*otpFn)
	if err != nil {
		log.Fatal(err)
	}

	jsonWebKey := jwt.GenerateJsonWebKey()

	fileServer := http.FileServer(http.Dir(*assetsDir))

	http.Handle("/", fileServer)

	// Define the HTTP handlers
	http.HandleFunc(
		"/token",
		jwt.GenerateTokenIssuer(
			jwt.WebKeyToSigner(jsonWebKey),
			authenticator,
		),
	)

	http.HandleFunc(
		"/jwks",
		jwt.GenerateJsonWebKeysRequestHandler(
			jwt.CreateKeySet(jwt.ExtractPublicJsonWebKey(jsonWebKey)),
		),
	)

	// Start the HTTP server
	log.Println("Athentication Server started on port 8008")
	log.Fatal(http.ListenAndServe(":8008", nil))
}
