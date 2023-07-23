package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gitamped/stem/data/nosql/dbschema"
	"github.com/gitamped/stem/database"
)

// Seed creates the schema in the database.
func Seed(ctx context.Context, cfg database.Config) error {

	b, _ := os.ReadFile("../../../testdata/seed.txt")
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
