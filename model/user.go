package model

type User struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
}

type UserWithCredentials struct {
	ID       string `bson:"_id,omitempty"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type UserResponse struct {
	User User `json:"user"`
	Errors []Error `json:"errors"`
}

// Should perhaps be moved to a helpers package?
func DowngradeToUser(userWithCredentials UserWithCredentials) User {
	return User{
		ID:    userWithCredentials.ID,
		Email: userWithCredentials.Email,
	}
}
