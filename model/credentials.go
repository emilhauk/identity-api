package model

type LoginRequestParams struct {
	RequestedUrl string `json:"requested_url,omitempty"`
	Credentials
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
