package jwt

import (
	"crypto/rand"
	"crypto/rsa"

	"gopkg.in/square/go-jose.v2"
)

func GenerateJsonWebKey() jose.JSONWebKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return jose.JSONWebKey{
		Algorithm: string(jose.PS512),
		Key:       key,
		KeyID:     KeyId,
	}
}

func ExtractPublicJsonWebKey(privateKey jose.JSONWebKey) jose.JSONWebKey {
	return jose.JSONWebKey{
		Algorithm: privateKey.Algorithm,
		Key:       &privateKey.Key.(*rsa.PrivateKey).PublicKey,
		KeyID:     privateKey.KeyID,
	}
}

func WebKeyToSigner(webKey jose.JSONWebKey) jose.Signer {
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(webKey.Algorithm),
			Key:       webKey.Key,
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		panic(err)
	}
	return signer
}

func CreateKeySet(webKeys ...jose.JSONWebKey) jose.JSONWebKeySet {
	return jose.JSONWebKeySet{
		Keys: webKeys,
	}
}
