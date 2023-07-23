package user

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService is an API for creating users for an app.
type UserService interface {
	// CreateUser create a user
	CreateUser(CreateUserRequest, server.GenericRequest) CreateUserResponse
	// UpdateUser updates a user
	UpdateUser(UpdateUserRequest, server.GenericRequest) UpdateUserResponse
	// DeleteUser deletes a user
	DeleteUser(DeleteUserRequest, server.GenericRequest) DeleteUserResponse
	// QueryUser retrieves a list of existing users
	QueryUser(QueryUserRequest, server.GenericRequest) QueryUserResponse
	// QueryByID gets the specified user by id
	QueryUserByID(QueryUserByIDRequest, server.GenericRequest) QueryUserByIDResponse
	// QueryByEmail gets the specified user by email
	QueryUserByEmail(QueryUserByEmailRequest, server.GenericRequest) QueryUserByEmailResponse
	// Authenticate finds a user by their email and verifies their password. On
	// success it returns a Claims User representing this user. The claims can be
	// used to generate a token for future authentication.
	Authenticate(AuthenticateRequest, server.GenericRequest) AuthenticateResponse
}

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) (User, error)
	Delete(ctx context.Context, email mail.Address) (User, error)
	QueryByID(ctx context.Context, id string) (User, error)
	QueryByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, usr UpdateUser) (User, error)
	Authenticate(ctx context.Context, email string, password string) (User, error)
}

// Required to register endpoints with the Server
type UserRpcService interface {
	UserService
	// Registers RPCService with Server
	Register(s *server.Server)
}

// Implements interface
type UserServicer struct {
	log    *zap.SugaredLogger
	storer Storer
	auth   auth.Auth
}

// Authenticate implements UserRpcService
func (u UserServicer) Authenticate(req AuthenticateRequest, gr server.GenericRequest) AuthenticateResponse {

	addr, err := mail.ParseAddress(req.Username)
	if err != nil {
		return AuthenticateResponse{Error: fmt.Errorf("invalid email format").Error()}

	}

	usr, err := u.storer.Authenticate(gr.Ctx, *&addr.Address, req.Password)
	if err != nil {
		return AuthenticateResponse{Error: err.Error()}
	}

	// flatten roles
	roles := make([]string, 0, len(usr.Roles))
	for _, value := range usr.Roles {
		roles = append(roles, value.name)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "bud project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: roles,
	}

	tkn, err := u.auth.GenerateToken(claims)
	if err != nil {
		return AuthenticateResponse{Error: fmt.Errorf("generatetoken: %w", err).Error()}

	}

	return AuthenticateResponse{Token: tkn}
}

// QueryUserByEmail implements UserRpcService
func (u UserServicer) QueryUserByEmail(req QueryUserByEmailRequest, gr server.GenericRequest) QueryUserByEmailResponse {
	usr, err := u.storer.QueryByEmail(gr.Ctx, req.Email)
	if err != nil {
		return QueryUserByEmailResponse{Error: err.Error()}
	}
	return QueryUserByEmailResponse{User: usr}
}

// QueryUserByID implements UserRpcService
func (u UserServicer) QueryUserByID(req QueryUserByIDRequest, gr server.GenericRequest) QueryUserByIDResponse {
	usr, err := u.storer.QueryByID(gr.Ctx, req.ID)
	if err != nil {
		return QueryUserByIDResponse{Error: err.Error()}
	}
	return QueryUserByIDResponse{User: usr}
}

// QueryUser implements UserRpcService
func (UserServicer) QueryUser(QueryUserRequest, server.GenericRequest) QueryUserResponse {
	panic("unimplemented")
}

// DeleteUser implements UserRpcService
func (u UserServicer) DeleteUser(req DeleteUserRequest, gr server.GenericRequest) DeleteUserResponse {
	du, err := u.storer.Delete(gr.Ctx, req.User.Email)
	if err != nil {
		return DeleteUserResponse{Error: err.Error()}
	}
	return DeleteUserResponse{User: du}
}

// CreateUser implements UserRpcService
func (u UserServicer) CreateUser(req CreateUserRequest, gr server.GenericRequest) CreateUserResponse {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return CreateUserResponse{Error: fmt.Errorf("generatefrompassword: %w", err).Error()}
	}
	usr := User{
		ID:           uuid.New(),
		Name:         req.NewUser.Name,
		Email:        req.NewUser.Email,
		PasswordHash: hash,
		Roles:        req.NewUser.Roles,
		Department:   req.NewUser.Department,
		Enabled:      true,
		DateCreated:  gr.Values.Now,
		DateUpdated:  gr.Values.Now,
	}
	result, err := u.storer.Create(gr.Ctx, usr)
	if err != nil {
		return CreateUserResponse{Error: err.Error()}
	}
	return CreateUserResponse{User: result}
}

// UpdateUser implements UserRpcService
func (u UserServicer) UpdateUser(req UpdateUserRequest, gr server.GenericRequest) UpdateUserResponse {
	uu, err := u.storer.Update(gr.Ctx, req.UpdateUser)
	if err != nil {
		return UpdateUserResponse{Error: err.Error()}
	}
	return UpdateUserResponse{User: uu}
}

// Register implements UserRpcService
func (us UserServicer) Register(s *server.Server) {
	s.Register("UserService", "CreateUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.CreateUserHandler})
	s.Register("UserService", "DeleteUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.DeleteUserHandler})
	s.Register("UserService", "QueryUserByID", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.QueryUserByIDHandler})
	s.Register("UserService", "QueryUserByEmail", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.QueryUserByEmailHandler})
	s.Register("UserService", "UpdateUser", server.RPCEndpoint{Roles: []string{auth.RoleAdmin}, Handler: us.UpdateUserHandler})
	s.Register("UserService", "Authenticate", server.RPCEndpoint{Roles: []string{}, Handler: us.AuthenticateHandler})
}

// Create new UserServicer
func NewUserServicer(log *zap.SugaredLogger, storer Storer, a auth.Auth) UserRpcService {
	return UserServicer{
		log:    log,
		storer: storer,
		auth:   a,
	}
}

// CreateUserRequest is the request object for UserService.CreateUser.
type CreateUserRequest struct {
	NewUser NewUser `json:"newUser"`
}

// CreateUserResponse is the response object containing a UserService.CreateUser.
type CreateUserResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

type UpdateUserRequest struct {
	UpdateUser UpdateUser `json:"user"`
}
type UpdateUserResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

// DeleteUserRequest is the request object for UserService.DeleteUser.
type DeleteUserRequest struct {
	User User `json:"user"`
}

// DeleteUserResponse is the response object for UserService.DeleteUser.
type DeleteUserResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

type QueryUserRequest struct{}
type QueryUserResponse struct{}

// QueryUserByIDRequest is the request object for UserService.QueryUserByID.
type QueryUserByIDRequest struct {
	ID string `json:"id"`
}

// QueryUserByIDResponse is the response object for UserService.QueryUserByID.
type QueryUserByIDResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

// QueryUserByEmailRequest is the request object for UserService.QueryUserByEmail.
type QueryUserByEmailRequest struct {
	Email string `json:"email"`
}

// QueryUserByEmailResponse is the response object for UserService.QueryUserByEmail.
type QueryUserByEmailResponse struct {
	User  User   `json:"user"`
	Error string `json:"error,omitempty"`
}

type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}
