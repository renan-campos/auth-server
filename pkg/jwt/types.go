package jwt

import "time"

type ClaimData struct {
	Subject  string
	Issuer   string
	Audience string
	Expiry   time.Time
	IssuedAt time.Time
}
