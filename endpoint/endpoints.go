package endpoint

import (
	"github.com/emilhauk/identity-api/config"
	"github.com/emilhauk/identity-api/store"
	"net/http"
)

type Endpoints struct {
	LoginHandler http.HandlerFunc
}

func NewEndpoints(store *store.MongoStore, config *config.Config) *Endpoints {

	endpoint := &Endpoints{
		LoginHandler: func(w http.ResponseWriter, r *http.Request) {
			LoginHandler(w, r, store, config.JwtSigningSecret)
		},
	}

	return endpoint
}
