package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/sirupsen/logrus"
)

type unverifiedClaims struct {
	model.RSAKeyIdentifier
	jwt.StandardClaims
}

func Keyfunc(keyStore *store.RSAKeyStore) jwt.Keyfunc {
	return func(token *jwt.Token) (key interface{}, err error) {
		var m unverifiedClaims
		_, _, err = new(jwt.Parser).ParseUnverified(token.Raw, &m)
		if err != nil {
			logrus.Errorln("KID missing from token", token)
			return nil, err
		}

		if keyPair, ok := keyStore.GetKeyPairById(m.KID); ok {
			return keyPair.Public, nil
		}
		return key, errors.New("No keypair matching kid: " + m.KID)
	}
}
