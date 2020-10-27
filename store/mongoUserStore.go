package store

import (
	"context"
	"github.com/emilhauk/identity-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type MongoUserStore struct {
	collection *mongo.Collection
}

func (s *MongoUserStore) FindByCredentials(credentials *model.Credentials) (user model.User, err error) {
	var userWithCredentials model.UserWithCredentials
	err = s.collection.FindOne(context.TODO(), bson.M{"username": credentials.Username}).Decode(&userWithCredentials)
	if err != nil {
		return user, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userWithCredentials.Password), []byte(credentials.Password))
	if err != nil {
		return user, err
	}
	return model.DowngradeToUser(userWithCredentials), nil
}

func (s *MongoUserStore) FindById(id string) (user model.User, err error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	err = s.collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
	return user, err
}
