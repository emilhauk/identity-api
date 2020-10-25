package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type RefreshToken struct {
	UserId string `bson:"user_id"`
	Token string `bson:"token"`
	Created primitive.Timestamp `bson:"created"`
	Expires primitive.Timestamp `bson:"expires"`
}
