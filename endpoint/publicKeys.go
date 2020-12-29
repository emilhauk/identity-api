package endpoint

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

func PublicKeyHandler(w http.ResponseWriter, r *http.Request, keyStore *store.RSAKeyStore) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	keys := map[string]string{}
	for keyId, keyPair := range keyStore.GetAllKeyPairs() {
		publicKey := &pem.Block{Type: "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(keyPair.Public)}
		keys[keyId] = string(pem.EncodeToMemory(publicKey))
	}

	response, err := json.Marshal(model.PublicKeysResponse{
		keys,
		[]model.Error{},
	})
	if err != nil {
		logrus.Errorln("Error marshalling PublicKeysResponse", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
