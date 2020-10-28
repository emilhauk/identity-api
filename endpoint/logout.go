package endpoint

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"log"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request, store *store.MongoStore, key []byte) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshTokenCookie, err := r.Cookie("refresh-token");
	if err != nil || len(refreshTokenCookie.Value) == 0 {
		log.Println("No refresh token present in cookie", r.Cookies())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var claims model.RefreshTokenClaims
	_, err = jwt.ParseWithClaims(refreshTokenCookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		log.Println("Invalid refresh token", refreshTokenCookie.Value, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = store.Token.DeleteByToken(claims.Token)
	if err != nil {
		log.Println("Failed deletion of refresh token")
	}

	// TODO: Perform invalidation of lingering JWTs?

	// TODO: Should also be able to delete specific keys

	w.WriteHeader(http.StatusNoContent)
}
