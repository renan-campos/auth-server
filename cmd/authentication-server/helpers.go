package main

import (
	"crypto/rand"
	"crypto/rsa"

	"gopkg.in/square/go-jose.v2"
)

func generateJsonWebKey() jose.JSONWebKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	kId := "my_rsa_key"

	return jose.JSONWebKey{
		Algorithm: string(jose.PS512),
		Key:       key,
		KeyID:     kId,
	}
}

func webKeyToSigner(webKey jose.JSONWebKey) jose.Signer {
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

func createKeySet(webKeys ...jose.JSONWebKey) jose.JSONWebKeySet {
	return jose.JSONWebKeySet{
		Keys: webKeys,
	}
}
