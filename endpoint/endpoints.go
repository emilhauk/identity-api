package endpoint

import (
	"github.com/emilhauk/identity-api/store"
	"net/http"
)

type Endpoints struct {
	LoginHandler  http.HandlerFunc
	JwtHandler    http.HandlerFunc
	LogoutHandler http.HandlerFunc
	WebHandler    http.HandlerFunc
	PublicKeyHandler http.HandlerFunc
}

func NewEndpoints(dbStore *store.MongoStore, keyStore *store.RSAKeyStore) *Endpoints {

	endpoint := &Endpoints{
		LoginHandler: func(w http.ResponseWriter, r *http.Request) {
			LoginHandler(w, r, dbStore, keyStore)
		},
		JwtHandler: func(w http.ResponseWriter, r *http.Request) {
			JwtHandler(w, r, dbStore, keyStore)
		},
		LogoutHandler: func(w http.ResponseWriter, r *http.Request) {
			LogoutHandler(w, r, dbStore, keyStore)
		},
		WebHandler: func(w http.ResponseWriter, r *http.Request) {
			WebHandler(w, r)
		},
		PublicKeyHandler: func(w http.ResponseWriter, r *http.Request) {
			PublicKeyHandler(w, r, keyStore)
		},
	}

	return endpoint
}
