package model

type LoginRequestParams struct {
	RequestedUrl string `json:"requested_url,omitempty"`
	Credentials
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
