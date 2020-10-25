package store

import "go.mongodb.org/mongo-driver/mongo"

type MongoStore struct {
	client *mongo.Client
	User MongoUserStore
	Token MongoTokenStore
}

func NewMongoStore(mongoClient *mongo.Client) MongoStore {
	database := mongoClient.Database("identity")

	userStore := MongoUserStore{
		collection: database.Collection("user"),
	}

	tokenStore := MongoTokenStore{
		collection: database.Collection("token"),
	}

	return MongoStore{
		client: mongoClient,
		User:   userStore,
		Token:  tokenStore,
	}
}
