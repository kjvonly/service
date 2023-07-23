package nosql

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/arangodb/go-driver"
	"github.com/gitamped/bud/services/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const collectionName = "users"

var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

type Store struct {
	db  driver.Database
	col driver.Collection
	log *zap.SugaredLogger
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db driver.Database) *Store {
	col, err := db.Collection(context.Background(), "users")
	if err != nil {
		log.Panicf("error accessing collection: %s", err)
	}
	return &Store{
		log: log,
		db:  db,
		col: col,
	}
}

// Delete deletes a user from the database
func (s *Store) Delete(ctx context.Context, email mail.Address) (user.User, error) {
	var result dbUser
	ctx = driver.WithReturnOld(ctx, &result)
	_, err := s.col.RemoveDocument(ctx, email.Address)
	return toCoreUser(result), err
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) (user.User, error) {
	var result dbUser
	ctx = driver.WithReturnNew(ctx, &result)
	_, err := s.col.CreateDocument(ctx, toDBUser(usr))
	return toCoreUser(result), err
}

// QueryById queries a user by id.
func (s *Store) QueryByID(ctx context.Context, id string) (user.User, error) {
	var result dbUser
	query := `FOR u IN @@coll
	FILTER u.user_id == @id
	LIMIT 1
	RETURN u`

	bindvars := map[string]interface{}{
		"@coll": collectionName,
		"id":    id,
	}

	c, err := s.db.Query(ctx, query, bindvars)
	defer c.Close()
	if err != nil {
		return user.User{}, err
	}
	_, err = c.ReadDocument(ctx, &result)
	return toCoreUser(result), err
}

// QueryById queries a user by email.
func (s *Store) QueryByEmail(ctx context.Context, email string) (user.User, error) {
	var result dbUser
	_, err := s.col.ReadDocument(ctx, email, &result)
	return toCoreUser(result), err
}

// Update updates a user by data.
func (s *Store) Update(ctx context.Context, updateUser user.UpdateUser) (user.User, error) {
	var result dbUser
	ctx = driver.WithReturnNew(ctx, &result)
	ctx = driver.WithKeepNull(ctx, false)
	_, err := s.col.UpdateDocument(ctx, updateUser.Email.Address, updateUser)
	return toCoreUser(result), err

}

func (s *Store) Authenticate(ctx context.Context, email string, password string) (user.User, error) {
	usr, err := s.QueryByEmail(ctx, email)
	if err != nil {
		return user.User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return user.User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
