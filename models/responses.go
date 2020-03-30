package models

// Response contains all the propeties of a response.
// It is used for error and standard success responses.
type Response struct {
	Ok      bool   `json:"ok"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
