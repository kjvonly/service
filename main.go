package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gitamped/bud/services/user"
	"github.com/gitamped/bud/services/user/stores/nosql"
	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/keystore"
	"github.com/gitamped/seed/mid"
	"github.com/gitamped/seed/server"
	"github.com/gitamped/stem/database"
	"go.uber.org/zap"
)

func main() {
	// New RPCServer
	s := server.NewServer(mid.CommonMiddleware)

	p, _ := filepath.Abs("./")
	fsPath := path.Join(p, "zarf", "keys")

	ks, _ := keystore.NewFS(os.DirFS(fsPath))

	a, _ := auth.New("54bb2165-71e1-41a6-af3e-7da4a0e1e2c1", ks)

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, err := database.Open(database.Config{
		User:       "root",
		Password:   "arangodb",
		Host:       fmt.Sprintf("http://%s", "127.0.0.1:49157"),
		Name:       "arangodb",
		DisableTLS: true,
	})

	if err != nil {
		sugar.Fatalf("Opening database connection: %v", err)
	}

	sugar.Info("Waiting for database to be ready ...")

	if err := database.StatusCheck(ctx, dbClient); err != nil {
		sugar.Fatalf("status check database: %v", err)
	}

	sugar.Info("Database ready")

	db, _ := dbClient.Database(ctx, "testcreateuser")
	userStorer := nosql.NewStore(sugar, db)

	// Register UserServicer
	gs := user.NewUserServicer(sugar, userStorer, *a)
	gs.Register(s)

	// Listen
	fmt.Println(`Listening on port 8080`)
	fmt.Println(`test cmd: curl -X POST  --data '{"username": "user@example.com", "password": "gophers"}' http://localhost:8080/v1/UserService.Authenticate`)
	http.Handle("/v1/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
