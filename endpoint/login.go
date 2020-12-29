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

func LoginHandler(w http.ResponseWriter, r *http.Request, dbStore *store.MongoStore, keyStore *store.RSAKeyStore) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginParams := model.LoginRequestParams{
		RequestedUrl: r.Form.Get("requested_url"),
		Credentials:  model.Credentials{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		},
	}

	user, err := dbStore.User.FindByCredentials(loginParams.Credentials)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, signedJwtToken, err := createRefreshToken(user, keyStore)
	if err != nil {
		logrus.Errorln("Unable to create signed refresh token.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = dbStore.Token.SaveToken(claims); err != nil {
		logrus.Errorln("Unable to save refresh token.", err)
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
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	logrus.Info(loginParams)

	if loginParams.RequestedUrl != "" {
		http.Redirect(w, r, loginParams.RequestedUrl, http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createRefreshToken(user model.User, keyStore *store.RSAKeyStore) (claims model.RefreshTokenClaims, signedJwt string, err error) {
	hash := sha256.New().Sum([]byte(util.RandomString(64)))
	refreshToken := hex.EncodeToString(hash[:])
	now := time.Now()
	expires := now.AddDate(0, 1, 0)

	claims = model.RefreshTokenClaims{
		Token: refreshToken,
		RSAKeyIdentifier: model.RSAKeyIdentifier{KID: keyStore.DefaultKeyId},
		StandardClaims: jwt.StandardClaims{
			Id: user.ID,
			ExpiresAt: expires.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	signedJwt, err = jwt.NewWithClaims(
		jwt.SigningMethodRS512,
		claims,
	).SignedString(keyStore.GetDefaultKeyPair().Private)
	return claims, signedJwt, err
}
