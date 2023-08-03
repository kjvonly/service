package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ardanlabs/conf/v3"

	"github.com/gitamped/seed/auth"
	"github.com/gitamped/seed/keystore"
	"github.com/gitamped/seed/mid"
	"github.com/gitamped/seed/server"
	"github.com/gitamped/stem/database"
	"github.com/kjvonly/service/services/bible"
	esStore "github.com/kjvonly/service/services/bible/stores/elasticsearch"
	"github.com/kjvonly/service/services/user"
	userStore "github.com/kjvonly/service/services/user/stores/nosql"
	"go.uber.org/zap"
)

var build = "develop"

type config struct {
	conf.Version
	Args       conf.Args
	ArangodbDB database.Config
}

func main() {

	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "kjvonly"
	_, err := conf.Parse(prefix, &cfg)

	if err != nil {
		log.Fatalf("failed to parse config")
	}

	esUrl := flag.String("esUrl", "http://127.0.0.1:9200", "elasticsearch url")
	flag.Parse()

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
		User:       cfg.ArangodbDB.User,
		Password:   cfg.ArangodbDB.Password,
		Host:       cfg.ArangodbDB.Host,
		Name:       cfg.ArangodbDB.Name,
		DisableTLS: cfg.ArangodbDB.DisableTLS,
	})

	if err != nil {
		sugar.Fatalf("Opening database connection: %v", err)
	}

	sugar.Info("Waiting for database to be ready ...")

	if err := database.StatusCheck(ctx, dbClient); err != nil {
		sugar.Fatalf("status check database: %v", err)
	}

	sugar.Info("Database ready")

	db, _ := dbClient.Database(ctx, "kjvonly")
	userStorer := userStore.NewStore(sugar, db)

	// Register UserServicer
	gs := user.NewUserServicer(sugar, userStorer, *a)
	gs.Register(s)

	// Register BibleSearchService
	ess := esStore.NewStore(sugar, *esUrl)
	bs := bible.NewBibleSearchServicer(sugar, ess, *a)
	bs.Register(s)

	// Listen
	fmt.Println(`Listening on port 8080`)
	fmt.Println(`test cmd: curl -X POST  --data '{"username": "user@example.com", "password": "gophers"}' http://localhost:8080/v1/UserService.Authenticate`)
	http.Handle("/v1/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
