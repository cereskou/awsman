package cognito

import "ditto.co.jp/awsman/model"

//Response -
type Response struct {
	Total  int64                 `json:"total"`  //user count
	Items  []*model.CognitoUser  `json:"items"`  //user list
	Groups []*model.CognitoGroup `json:"groups"` //group list
}

//ChangePasswordRequest - password
type ChangePasswordRequest struct {
	UUID      string `json:"uuid"`
	Password  string `json:"password"`
	Permanent string `json:"permanent"`
}
