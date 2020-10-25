package store

import (
	"context"
	"github.com/emilhauk/identity-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type MongoUserStore struct {
	collection *mongo.Collection
}

func (s *MongoUserStore) FindByCredentials(credentials *model.Credentials) (model.User, error) {
	var user model.User // only used to return empty user
	var userWithCredentials model.UserWithCredentials
	err := s.collection.FindOne(context.TODO(), bson.D{{"username", credentials.Username}}).Decode(&userWithCredentials)
	if err != nil {
		return user, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userWithCredentials.Password), []byte(credentials.Password))
	if err != nil {
		return user, err
	}
	return userWithCredentials.User, err
}
