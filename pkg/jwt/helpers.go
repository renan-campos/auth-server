package jwt

import (
	"crypto/rand"
	"crypto/rsa"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
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

func NewClaims(claimData ClaimData) jwt.Claims {
	return jwt.Claims{
		Subject:  claimData.Subject,
		Issuer:   claimData.Issuer,
		Audience: jwt.Audience{claimData.Audience},
		Expiry:   jwt.NewNumericDate(claimData.Expiry),
		IssuedAt: jwt.NewNumericDate(claimData.IssuedAt),
	}
}

func CreateToken(signer jose.Signer, claims jwt.Claims) (string, error) {
	return jwt.Signed(signer).Claims(claims).CompactSerialize()
}
