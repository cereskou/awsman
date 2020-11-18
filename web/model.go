package web

//User -
type User struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Roles        []string `json:"roles"`
}

//Token -
type Token struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Expires      int64  `json:"expires"`
}

//Login -
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
