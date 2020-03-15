package models

type Setup struct {
	Organization Organization `json:"organization"`
	User         User         `json:"user"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	User
	JWT string `json:"jwt"`
}
