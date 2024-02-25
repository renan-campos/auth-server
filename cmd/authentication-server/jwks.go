/*
JWKS Endpoint:

  - This endpoint is responsible for serving the JSON Web Key Set (JWKS), which
    contains the public keys that clients can use to verify JWT signatures.

  - Conventionally, this endpoint is often named /jwks or /oauth/jwks.

  - It usually serves JWKS documents in JSON format containing one or more JSON Web
    Keys (JWKs). Clients can retrieve the JWKS from this endpoint to dynamically
    obtain the public keys needed for JWT verification.

  - This endpoint is often associated with OAuth 2.0 and OpenID Connect
    implementations.
*/
package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/square/go-jose.v2"
)

func generateJsonWebKeysRequestHandler(
	webKeySet jose.JSONWebKeySet,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Process the GET request
		// You can add your logic here to handle the GET request
		marshalledWebKeySet, err := json.Marshal(webKeySet)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		// Send a success response
		w.WriteHeader(http.StatusOK)
		w.Write(marshalledWebKeySet)
	}
}
