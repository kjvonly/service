// This program performs administrative tasks for the garage sale service.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/gitamped/stem/database"

	"github.com/kjvonly/service/tooling/services/kjvonly-admin/commands"
	"go.uber.org/zap"
)

var build = "develop"

type config struct {
	conf.Version
	Args       conf.Args
	ArangodbDB database.Config
	Seed       struct {
		Path string `conf:"default:testdata/seed.txt"`
	}
	Migrate struct {
		Path string `conf:"default:testdata/collections.txt"`
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	if err := run(log); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("msg", err)
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "kjvonly"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		log.Info(context.Background(), "startup", "config", out)

		return fmt.Errorf("parsing config: %w", err)
	}

	return processCommands(cfg.Args, log, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, log *zap.SugaredLogger, cfg config) error {
	ctx := context.TODO()

	switch args.Num(0) {
	case "migrate":
		if err := commands.Migrate(ctx, cfg.ArangodbDB, cfg.Migrate.Path); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}

	case "seed":
		if err := commands.Seed(ctx, cfg.ArangodbDB, cfg.Seed.Path); err != nil {
			return fmt.Errorf("seeding database: %w", err)
		}

	default:
		fmt.Println("migrate:    create the schema in the database")
		fmt.Println("seed:       add data to the database")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
