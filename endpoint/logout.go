package endpoint

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/emilhauk/identity-api/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request, store *store.MongoStore, key []byte) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshTokenCookie, err := r.Cookie("refresh-token")
	if err != nil || len(refreshTokenCookie.Value) == 0 {
		logrus.Debugln("refresh-token", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var claims model.RefreshTokenClaims
	_, err = jwt.ParseWithClaims(refreshTokenCookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		logrus.Debugln("Invalid refresh token", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = store.Token.DeleteByToken(claims.Token)
	if err != nil {
		logrus.Errorln("Failed deletion of refresh token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Perform invalidation of lingering JWTs?
	// TODO: Should also be able to delete specific keys

	util.DeleteCookie(refreshTokenCookie, w)
	w.WriteHeader(http.StatusOK)
}
