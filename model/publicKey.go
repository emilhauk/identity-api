package model

import "crypto/rsa"

type PublicKeysResponse struct {
	PublicKeys map[string]string `json:"public_keys"`
	Errors []Error `json:"errors"`
}

type RSAKeyPair struct {
	Public *rsa.PublicKey
	Private *rsa.PrivateKey
}
