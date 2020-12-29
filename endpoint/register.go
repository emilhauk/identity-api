package endpoint

import (
	"encoding/json"
	"github.com/emilhauk/identity-api/model"
	"github.com/emilhauk/identity-api/store"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, store *store.MongoStore) {
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

	registerParams := model.RegisterRequestParams{
		r.Form.Get("name"),
		r.Form.Get("email"),
		r.Form.Get("password"),
	}
	if registerParams.Name == "" {
		registerParams.Name = strings.Split(registerParams.Email, "@")[0]
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerParams.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	registerParams.Password = string(hashedPassword)

	user, err := store.User.Create(registerParams)
	if err != nil {
		logrus.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logrus.Infof("Created: (%s)", user)
	response, err := json.Marshal(model.UserResponse{
		User:   user,
		Errors: []model.Error{},
	})
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
