package endpoint

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func JwtHandler(w http.ResponseWriter, r *http.Request, store *store.MongoStore, key []byte) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:9000")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshTokenCookie, err := r.Cookie("refresh-token")
	if err != nil || len(refreshTokenCookie.Value) == 0 {
		logrus.Println("No refresh token present in cookie", r.Cookies())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var claims model.RefreshTokenClaims
	_, err = jwt.ParseWithClaims(refreshTokenCookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
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

	user, err := store.User.FindById(dbToken.UserId)
	if err != nil {
		logrus.Println("User mentioned by db not in database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	userClaims := model.UserTokenClaims{
		Id:   user.ID,
		Name: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Duration(30 * time.Minute)).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	signedJwt, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		userClaims,
	).SignedString(key)
	w.Header().Add("Authorization", signedJwt)
	w.WriteHeader(204)
}
