package endpoint

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/emilhauk/identity-api/util"
	"github.com/sirupsen/logrus"
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

	claims, signedJwtToken, err := createRefreshToken(user, key)
	if err != nil {
		logrus.Println("Unable to create signed refresh token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = store.Token.SaveToken(user.ID, claims); err != nil {
		logrus.Println("Unable to save refresh token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh-token",
		Value:    signedJwtToken,
		Path:     "/",
		Domain:   r.Header.Get("host"),
		Expires:  time.Unix(claims.ExpiresAt, 0),
		Secure:   r.Host == "localhost",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(204)
}

func createRefreshToken(user model.User, key []byte) (claims model.RefreshTokenClaims, signedJwt string, err error) {
	hash := sha256.New().Sum([]byte(util.RandomString(64)))
	refreshToken := hex.EncodeToString(hash[:])
	now := time.Now()
	expires := now.AddDate(0, 1, 0)

	claims = model.RefreshTokenClaims{
		Token: refreshToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	signedJwt, err = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(key)
	return claims, signedJwt, err
}
