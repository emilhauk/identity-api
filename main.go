package main

import (
	"context"
	"github.com/emilhauk/identity-api/config"
	"github.com/emilhauk/identity-api/endpoint"
	"github.com/emilhauk/identity-api/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func setupRoutes(endpoints *endpoint.Endpoints) {
	http.HandleFunc("/login", endpoints.LoginHandler)
	http.HandleFunc("/jwt", endpoints.JwtHandler)
	http.HandleFunc("/logout", endpoints.LogoutHandler)
}

func main() {
	c := config.NewConfig()

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(c.MongoDBUrl))
	if err != nil {
		log.Panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)
	mongoStore := store.NewMongoStore(mongoClient)

	endpoints := endpoint.NewEndpoints(&mongoStore, c)

	setupRoutes(endpoints)
	log.Fatal(http.ListenAndServe(c.Host, nil))
}
