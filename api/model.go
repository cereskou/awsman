package api

//Response -
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// HTTPError represents an error that occurred while handling a request.
// fork echo.HTTPError for swagger
type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}
