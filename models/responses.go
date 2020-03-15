package models

type ErrorResponse struct {
	Ok      bool   `json:"ok"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
