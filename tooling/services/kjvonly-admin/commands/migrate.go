package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"git.launchpad.net/~man4christ/+git/stem/data/nosql/dbschema"
	"git.launchpad.net/~man4christ/+git/stem/database"
)

var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(ctx context.Context, cfg database.Config, migratePath string) error {

	b, _ := os.ReadFile(migratePath)
	cols := strings.Split(string(b), "\n")

	dbClient, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("Opening database connection: %v", err)
	}

	log.Println("Waiting for database to be ready ...")

	if err := database.StatusCheck(ctx, dbClient); err != nil {
		return fmt.Errorf("status check database: %v", err)
	}

	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.CreateDatabase(ctx, dbClient, cfg)

	if err != nil {
		return fmt.Errorf("error creating database")
	}

	if err := dbschema.Migrate(ctx, db, cols, []string{}); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}
