package model

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
}

type UserWithCredentials struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// Should perhaps be moved to a helpers package?
func DowngradeToUser(userWithCredentials UserWithCredentials) User {
	return User{
		ID:       userWithCredentials.ID,
		Username: userWithCredentials.Username,
	}
}
