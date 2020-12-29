package endpoint

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/emilhauk/identity-api/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func JwtHandler(w http.ResponseWriter, r *http.Request, store *store.MongoStore, keyStore *store.RSAKeyStore) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:9000")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	refreshTokenCookie, err := r.Cookie("refresh-token")
	if err != nil || len(refreshTokenCookie.Value) == 0 {
		logrus.Println("No refresh token present in cookie", r.Cookies())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var claims model.RefreshTokenClaims
	_, err = jwt.ParseWithClaims(refreshTokenCookie.Value, &claims, util.Keyfunc(keyStore))
	if err != nil {
		logrus.Println("Invalid refresh token", refreshTokenCookie.Value, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbToken, err := store.Token.FindByToken(claims.Token)
	if err != nil {
		logrus.Println("Refresh token not found on server")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := store.User.FindById(dbToken.Id)
	if err != nil {
		logrus.Println("User mentioned by db not in database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	userClaims := model.UserTokenClaims{
		UserId:           user.ID,
		Email: 			  user.Email,
		RSAKeyIdentifier: model.RSAKeyIdentifier{keyStore.DefaultKeyId},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(30 * time.Minute).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	signedJwt, err := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		userClaims,
	).SignedString(keyStore.GetDefaultKeyPair().Private)
	w.Header().Add("Authorization", signedJwt)
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
	w.WriteHeader(201)
}
