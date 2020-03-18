package models

type Response struct {
	Ok      bool   `json:"ok"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
