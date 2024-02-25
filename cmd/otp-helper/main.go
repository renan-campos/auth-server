package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/xlzd/gotp"
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
	log.Print("Displaying URI as QR code to register")
	uri := authenticator.totp.ProvisioningUri("test-account", "test-issuer")
	// Generate a 'dense' qrcode with the 'Low' level error correction and write it to Stdout
	qrterminal.Generate(uri, qrterminal.L, os.Stdout)

	var valid bool
	var otp string
	for !valid {
		log.Print("Enter an otp to try to validate: ")
		fmt.Scan(&otp)
		valid = authenticator.totp.VerifyTime(otp, time.Now())
		if !valid {
			log.Printf("otp %q is not valid, please try again", otp)
		}
	}
}

type otpAuthenticator struct {
	totp *gotp.TOTP
}

func NewOtpAuthenticator(secretFileName string) (*otpAuthenticator, error) {
	fp, err := os.Open(secretFileName)
	if err != nil {
		return nil, err
	}

	out, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return &otpAuthenticator{
		totp: gotp.NewDefaultTOTP(string(out[:32])),
	}, nil
}
