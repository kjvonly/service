package nosql

import (
	"database/sql"
	"net/mail"
	"time"

	"github.com/gitamped/bud/services/user"
	"github.com/google/uuid"
)

// dbUser represent the structure we need for moving data
// between the app and the database.
type dbUser struct {
	ID           uuid.UUID      `json:"user_id"`
	Name         string         `json:"name"`
	Email        string         `json:"_key"`
	Roles        []string       `json:"roles"`
	PasswordHash []byte         `json:"password_hash"`
	Enabled      bool           `json:"enabled"`
	Department   sql.NullString `json:"department"`
	DateCreated  time.Time      `json:"date_created"`
	DateUpdated  time.Time      `json:"date_updated"`
}

func toDBUser(usr user.User) dbUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return dbUser{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		DateCreated:  usr.DateCreated.UTC(),
		DateUpdated:  usr.DateUpdated.UTC(),
	}
}

func toCoreUser(dbUsr dbUser) user.User {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roles := make([]user.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		roles[i] = user.MustParseRole(value)
	}

	usr := user.User{
		ID:           dbUsr.ID,
		Name:         dbUsr.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		DateCreated:  dbUsr.DateCreated.In(time.Local),
		DateUpdated:  dbUsr.DateUpdated.In(time.Local),
	}

	return usr
}

func toCoreUserSlice(dbUsers []dbUser) []user.User {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		usrs[i] = toCoreUser(dbUsr)
	}
	return usrs
}
