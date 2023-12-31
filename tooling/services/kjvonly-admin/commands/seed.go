package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"git.launchpad.net/~man4christ/+git/stem/data/nosql/dbschema"
	"git.launchpad.net/~man4christ/+git/stem/database"
)

// Seed populates the schema in the database.
func Seed(ctx context.Context, cfg database.Config, seedPath string) error {

	b, _ := os.ReadFile(seedPath)
	seed := string(b)

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

	if err := dbschema.Seed(ctx, db, seed); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("seed complete")
	return nil
}
