package endpoint

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/emilhauk/identity-api/util"
	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, store *store.MongoStore, key []byte) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var credentials *model.Credentials
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := store.User.FindByCredentials(credentials)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshToken := string(sha256.New().Sum([]byte(util.RandomString(64))))
	signedToken, expires, err := createRefreshToken(refreshToken, user, key)
	if err != nil {
		log.Println("Unable to create signed refresh token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = store.Token.SaveToken(user.ID, refreshToken, expires); err != nil {
		log.Println("Unable to save refresh token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:       "refresh-token",
		Value:      signedToken,
		Path:       "/identity",
		Domain:     r.Header.Get("host"),
		Expires:    expires,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(204)
}

func createRefreshToken(refreshToken string, user model.User, key []byte) (signedToken string, expires time.Time, err error) {
	expires = time.Now().AddDate(0, 1, 0)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expires.Unix(),
		"uid": user.ID,
		"token": refreshToken,
	})
	signedToken, err = token.SignedString(key)
	return signedToken, expires, err
}
