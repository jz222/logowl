package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type TeamMember struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
}

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
	IsOrganizationOwner bool               `json:"isOrganizationOwner" bson:"isOrganizationOwner"`
	Services            []Service          `json:"services" bson:"services"`
	Team                []TeamMember       `json:"team" bson:"team"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (u *User) Validate() bool {
	if u.FirstName == "" || u.LastName == "" || u.Email == "" || u.Role == "" || len(u.Password) < 12 {
		return false
	}

	return true
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
