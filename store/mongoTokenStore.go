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

func (s *MongoTokenStore) FindByToken(token string) (refreshToken model.RefreshToken, err error) {
	err = s.collection.FindOne(context.TODO(), bson.D{{"refreshtokenclaims.token", token}}).Decode(&refreshToken)
	return refreshToken, err
}

func (s *MongoTokenStore) SaveToken(id string, claims model.RefreshTokenClaims) error {
	refreshToken := model.RefreshToken{
		UserId:             id,
		RefreshTokenClaims: claims,
	}
	_, err := s.collection.InsertOne(context.TODO(), refreshToken)
	return err
}

func (s *MongoTokenStore) DeleteByToken(token string) error {
	res, err := s.collection.DeleteOne(context.TODO(), bson.M{"refreshtokenclaims.token": token})
	if err != nil && res.DeletedCount == 0 {
		logrus.Warn("Deleted 0 rows...")
	}
	return err
}
