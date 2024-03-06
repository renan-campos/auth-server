package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"os"

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

func KeyFromJson(jsonFileName string) (jose.JSONWebKey, error) {
	var out jose.JSONWebKey
	fp, err := os.Open(jsonFileName)
	if err != nil {
		return out, err
	}
	defer fp.Close()
	data, err := io.ReadAll(fp)
	if err != nil {
		return out, err
	}
	json.Unmarshal(data, &out)
	return out, nil
}

func KeyFromPemFile(pemFileName string) (jose.JSONWebKey, error) {
	var out jose.JSONWebKey

	fp, err := os.Open(pemFileName)
	if err != nil {
		return out, err
	}
	pemData, err := io.ReadAll(fp)
	if err != nil {
		return out, err
	}
	pemBlock, _ := pem.Decode(pemData)
	if pemBlock == nil || pemBlock.Type != "RSA PRIVATE KEY" {
		return out, err
	}

	key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return out, err
	}
	return jose.JSONWebKey{
		Algorithm: string(jose.PS256),
		Key:       key,
		KeyID:     "NPal4M9u9jYj3arqXLG6X1xEjz-BRU3qglt5taMlbRY",
	}, nil
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
