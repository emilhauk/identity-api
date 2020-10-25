package store

import (
	"context"
	"github.com/emilhauk/identity-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoTokenStore struct {
	collection *mongo.Collection
}

func (s *MongoTokenStore) FindByToken(token string) (model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := s.collection.FindOne(context.TODO(), bson.D{{"token", token}}).Decode(&refreshToken)
	return refreshToken, err
}

func (s *MongoTokenStore) SaveToken(id string, token string, expires time.Time) error {
	refreshToken := model.RefreshToken{
		UserId:  id,
		Token:   token,
		Expires: primitive.Timestamp{T: uint32(expires.Unix())},
		Created: primitive.Timestamp{T: uint32(time.Now().Unix())},
	}
	_, err := s.collection.InsertOne(context.TODO(), refreshToken)
	return err
}
