package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Email        mail.Address `json:"email"`
	Roles        []Role       `json:"roles"`
	PasswordHash []byte       `json:"password_hash"`
	Department   string       `json:"department"`
	Enabled      bool         `json:"enabled"`
	DateCreated  time.Time    `json:"date_created"`
	DateUpdated  time.Time    `json:"date_updated"`
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string       `json:"name"`
	Email           mail.Address `json:"email"`
	Roles           []Role       `json:"roles"`
	Department      string       `json:"department"`
	Password        string       `json:"password"`
	PasswordConfirm string       `json:"password_confirm"`
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name            *string       `json:"name"`
	Email           *mail.Address `json:"email"`
	Roles           []Role        `json:"roles"`
	Department      *string       `json:"department"`
	Password        *string       `json:"passowrd"`
	PasswordConfirm *string       `json:"password_confirm"`
	Enabled         *bool         `json:"enabled"`
}
