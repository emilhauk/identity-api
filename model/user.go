package model

type User struct {
	ID string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
}

type UserWithCredentials struct {
	User
	Password string `bson:"password"`
}
