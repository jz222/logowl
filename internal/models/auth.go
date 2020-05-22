package models

import "time"

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
	JWT            string `json:"jwt,omitempty"`
	AccessPass     string `json:"accessPass,omitempty"`
	ExpirationTime int64  `json:"expirationTime"`
}

// PasswordResetToken contains the password reset token and the associated email address.
type PasswordResetToken struct {
	Email     string    `json:"email" bson:"email"`
	Token     string    `json:"token" bson:"token"`
	Used      bool      `json:"used" bson:"used"`
	ExpiresAt int64     `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type PasswordResetBody struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
}
