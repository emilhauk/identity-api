package main

import (
	"context"
	"github.com/emilhauk/identity-api/config"
	"github.com/emilhauk/identity-api/endpoint"
	"github.com/emilhauk/identity-api/store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func setupRoutes(endpoints *endpoint.Endpoints) {
	http.HandleFunc("/register", endpoints.RegisterHandler)
	http.HandleFunc("/login", endpoints.LoginHandler)
	http.HandleFunc("/jwt", endpoints.JwtHandler)
	http.HandleFunc("/logout", endpoints.LogoutHandler)
	http.HandleFunc("/keys", endpoints.PublicKeyHandler)
	http.HandleFunc("/", endpoints.WebHandler)
}

func main() {
	c := config.NewConfig()

	logrus.SetLevel(c.LogLevel)

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(c.MongoDBUrl))
	if err != nil {
		logrus.Panicln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer mongoClient.Disconnect(ctx)
	mongoStore := store.NewMongoStore(mongoClient)
	keyStore := store.NewRSAKeyStore(c.KeyStorePath, c.DefaultKeyId)
	endpoints := endpoint.NewEndpoints(&mongoStore, &keyStore)

	setupRoutes(endpoints)

	logrus.Infof("Identity-service listening on %s", c.Host)
	logrus.Fatalln(http.ListenAndServe(c.Host, nil))
}
