package user_test

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gitamped/bud/services/user"
	"github.com/gitamped/bud/services/user/stores/nosql"
	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/server"
	"github.com/gitamped/seed/values"
	"github.com/gitamped/stem/data/nosql/dbtest"
	"github.com/gitamped/stem/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func Test_User(t *testing.T) {
	b, _ := os.ReadFile("../../testdata/collections.txt")
	cols := strings.Split(string(b), "\n")
	b, _ = os.ReadFile("../../testdata/seed.txt")
	seed := string(b)

	d := dbtest.Data{
		CollectionData: cols,
		SeedAql:        seed,
	}

	test := dbtest.NewIntegration(t, c, "testuser", d)
	log := test.Log
	db := test.DB
	teardown := test.Teardown
	authSvc := test.Auth
	t.Cleanup(teardown)
	storer := nosql.NewStore(log, db)

	core := user.NewUserServicer(log, storer, *authSvc)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
			email, err := mail.ParseAddress("user@example.com")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to parse email: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse email.", dbtest.Success, testID)

			nu := user.CreateUserRequest{}
			nu.NewUser.Name = "John Doe"
			nu.NewUser.Email = *email
			nu.NewUser.Roles = []user.Role{user.RoleAdmin}
			nu.NewUser.Password = "gophers"
			nu.NewUser.PasswordConfirm = "gophers"

			cuUsr := core.CreateUser(nu, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})
			if cuUsr.User.Name != "John Doe" {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user %+v : got %+v.", dbtest.Failed, testID, nu, cuUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", dbtest.Success, testID)

			// query user by id
			qu := user.QueryUserByIDRequest{cuUsr.User.ID.String()}
			quUsr := core.QueryUserByID(qu, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if quUsr.User.ID != cuUsr.User.ID && quUsr.User.Email != cuUsr.User.Email {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query user by id %+v : got %+v.", dbtest.Failed, testID, qu, quUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to query user by id.", dbtest.Success, testID)

			// query user by email
			que := user.QueryUserByEmailRequest{Email: cuUsr.User.Email.Address}
			queUsr := core.QueryUserByEmail(que, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if queUsr.User.ID != cuUsr.User.ID && queUsr.User.Email != cuUsr.User.Email {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query user by email %+v : got %+v.", dbtest.Failed, testID, que, queUsr)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to query user by email.", dbtest.Success, testID)

			// update user
			var updateName string = "updated user name"
			uusr := user.UpdateUser{
				Email: &cuUsr.User.Email,
				Name:  &updateName,
			}
			uu := user.UpdateUserRequest{uusr}
			uuUsr := core.UpdateUser(uu, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{Roles: []string{auth.RoleAdmin}},
				Values: &values.Values{Now: now},
			})

			if uuUsr.User.Email.Address != uu.UpdateUser.Email.Address || uuUsr.User.ID != cuUsr.User.ID {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user %+v : got %+v.", dbtest.Failed, testID, uu, uuUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", dbtest.Success, testID)

			// authenticat user
			au := user.AuthenticateRequest{
				Username: email.Address,
				Password: "gophers",
			}

			auUsr := core.Authenticate(au, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if len(auUsr.Token) == 0 {
				t.Fatalf("\t%s\tTest %d:\tShould be able to authenticate user %+v : got %+v.", dbtest.Failed, testID, au, auUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to authenticate user.", dbtest.Success, testID)

			// authenticat user
			auf := user.AuthenticateRequest{
				Username: email.Address,
				Password: "wrong password",
			}

			aufUsr := core.Authenticate(auf, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if aufUsr.Error != "comparehashandpassword: authentication failed" {
				t.Fatalf("\t%s\tTest %d:\tShould be able to forbid failed authenticated user %+v : got %+v.", dbtest.Failed, testID, auf, aufUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to forbid failed authenticated user.", dbtest.Success, testID)

			// delete user
			du := user.DeleteUserRequest{cuUsr.User}
			duUsr := core.DeleteUser(du, server.GenericRequest{
				Ctx:    ctx,
				Claims: auth.Claims{},
				Values: &values.Values{Now: now},
			})

			if duUsr.User.ID != cuUsr.User.ID {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user %+v : got %+v.", dbtest.Failed, testID, du, duUsr)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", dbtest.Success, testID)

		}
	}
}
