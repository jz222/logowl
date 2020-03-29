package models

type Setup struct {
	Organization Organization `json:"organization"`
	User         User         `json:"user"`
}

type Credentials struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode,omitempty"`
}

type SignInResponse struct {
	User
	JWT            string `json:"jwt"`
	ExpirationTime int64  `json:"expirationTime"`
}
