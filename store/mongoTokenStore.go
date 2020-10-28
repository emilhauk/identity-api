package store

import (
	"context"
	"github.com/emilhauk/identity-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
	res, err := s.collection.DeleteOne(context.TODO(), bson.M{"token": token})
	if res.DeletedCount == 0 {
		log.Println("Deleted 0 rows...")
	}
	return err
}
