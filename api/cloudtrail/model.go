package cloudtrail

//Resource -
type Resource struct {
	Name string `json:"resourceName"`
	Type string `json:"resourceType"`
}

//Event -
type Event struct {
	EventID         string      `json:"id"`
	AccessKeyID     string      `json:"accessKey"`
	CloudTrailEvent string      `json:"event"`
	EventName       string      `json:"name"`
	EventSource     string      `json:"source"`
	EventTime       string      `json:"eventtime"`
	Username        string      `json:"username"`
	Resources       []*Resource `json:"resources"`
}

//Response -
type Response struct {
	Events    []*Event `json:"events"` //event list
	NextToken string   `json:"token"`  //next token
}
