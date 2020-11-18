package code

//Repository -
type Repository struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Branchs []string `json:"branchs"`
}

//Response -
type Response struct {
	Repositories []*Repository `json:"items"` //event list
	NextToken    string        `json:"token"` //next token
}
