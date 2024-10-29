package otp

import (
	"encoding/base32"
	"io"
	"net/http"
	"os"
	"time"
	"fmt"

	"github.com/xlzd/gotp"
)

type authenticator struct {
	totp *gotp.TOTP
}

func NewAuthenticator(secretFileName string) (Authenticator, error) {
	secret, err := extractSecretFromFile(secretFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to extract secret from file: %v", err)
	}

	if err := validateSecret(secret); err != nil {
		return nil, fmt.Errorf("invalid secret: %v", err)
	}

	return &authenticator{
		totp: gotp.NewDefaultTOTP(secret),
	}, nil
}

func extractSecretFromFile(secretFileName string) (string, error) {
	fp, err := os.Open(secretFileName)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	out, err := io.ReadAll(fp)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func validateSecret(secret string) error {
	// The current implementation expects the secret to be a valid base32 encoded string.
	if _, err := base32.StdEncoding.DecodeString(secret); err != nil {
		return err
	}
	return nil
}

func (a *authenticator) VerifyHttpRequest(r *http.Request) (valid bool, err error) {
	otp, err := io.ReadAll(r.Body)
	if err != nil {
		return valid, err
	}
	return a.totp.VerifyTime(string(otp), time.Now()), nil
}

func (a *authenticator) CurrentToken() string {
	return a.totp.Now()
}
