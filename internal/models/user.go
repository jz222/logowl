package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// TeamMember is a leaner version of an user.
type TeamMember struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	FirstName           string             `json:"firstName" bson:"firstName"`
	LastName            string             `json:"lastName" bson:"lastName"`
	Email               string             `json:"email" bson:"email"`
	Role                string             `json:"role" bson:"role"`
	IsVerified          bool               `json:"isVerified" bson:"isVerified"`
	InviteCode          string             `json:"inviteCode,omitempty" bson:"inviteCode,omitempty"`
	IsOrganizationOwner bool               `json:"isOrganizationOwner,omitempty" bson:"isOrganizationOwner,omitempty"`
}

// User contains information about the user and their
// organization and all services they have access to.
type User struct {
	ID                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName           string             `json:"firstName" bson:"firstName"`
	LastName            string             `json:"lastName" bson:"lastName"`
	Email               string             `json:"email" bson:"email"`
	Password            string             `json:"password,omitempty" bson:"password"`
	Role                string             `json:"role" bson:"role"`
	OrganizationID      primitive.ObjectID `json:"organizationId" bson:"organizationId"`
	Organization        *Organization      `json:"organization,omitempty" bson:"organization,omitempty"`
	IsVerified          bool               `json:"-" bson:"isVerified"`
	InviteCode          string             `json:"inviteCode" bson:"inviteCode,omitempty"`
	IsOrganizationOwner bool               `json:"isOrganizationOwner" bson:"isOrganizationOwner"`
	Services            []Service          `json:"services" bson:"services"`
	Team                []TeamMember       `json:"team,omitempty" bson:"team,omitempty"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Validate validates the data of an user.
func (u *User) Validate() bool {
	if u.FirstName == "" || u.LastName == "" || u.Email == "" || u.Role == "" || len(u.Password) < 8 {
		return false
	}

	return true
}

// VerifyPassword verifies a given password against the user's hashed password.
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}

// IsAdmin determines if the user is an admin.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsOwnerOfOrganization determines if the user is the owner of their organization.
func (u *User) IsOwnerOfOrganization() bool {
	return u.IsOrganizationOwner
}
