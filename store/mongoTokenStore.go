package store

import (
	"context"
	"github.com/emilhauk/identity-api/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTokenStore struct {
	collection *mongo.Collection
}

func (s *MongoTokenStore) FindByToken(token string) (claims model.RefreshTokenClaims, err error) {
	err = s.collection.FindOne(context.TODO(), bson.D{{"token", token}}).Decode(&claims)
	return claims, err
}

func (s *MongoTokenStore) SaveToken(claims model.RefreshTokenClaims) error {
	_, err := s.collection.InsertOne(context.TODO(), claims)
	return err
}

func (s *MongoTokenStore) DeleteByToken(token string) error {
	res, err := s.collection.DeleteOne(context.TODO(), bson.M{"token": token})
	if err != nil && res.DeletedCount == 0 {
		logrus.Warn("Deleted 0 rows...")
	}
	return err
}
