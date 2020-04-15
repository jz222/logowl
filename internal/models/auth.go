package models

// Setup contains all properties related to the organization setup.
type Setup struct {
	Organization Organization `json:"organization"`
	User         User         `json:"user"`
}

// Credentials contains the payload for signing in and signing up.
type Credentials struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode,omitempty"`
}

// SignInResponse contains the response for a successfull sign in.
type SignInResponse struct {
	User
	JWT            string `json:"jwt"`
	ExpirationTime int64  `json:"expirationTime"`
}
