package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/miodzie/gobroke/app"
	"github.com/miodzie/gobroke/config"
	"github.com/miodzie/gobroke/notifying"
	"github.com/miodzie/gobroke/notifying/email"
	"github.com/miodzie/gobroke/pricing"
	"github.com/miodzie/gobroke/storage/memory"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var (
		port       = flags.Int("port", 8001, "port to serve api")
		debug      = flags.Bool("debug", false, "sets log level to debug")
		configPath = flags.String("config", "~/.config/gobroke/config.yaml", "specify config")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	initLogs(*debug)
	cfg, err := initConfig(*configPath)
	if err != nil {
		return err
	}

	c, err := initChecker(cfg)
	if err != nil {
		return err
	}

	// TODO: refactor later.
	db, err := sql.Open("sqlite3", "~/.condig/gobroke/database.sqlite")
	if err != nil {
		return err
	}

	trigRepo := memory.NewTriggerRepo()
	t := &notifying.Trigger{ID: 634, Symbol: "BTC"}
	trigRepo.Save(t)
	srv, err := app.NewServer(db, c, trigRepo)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	fmt.Printf("listening on: %d\n", *port)

	return http.ListenAndServe(addr, srv)
}

func initLogs(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

  log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// TODO: Add check for config debug too. Flag should always override.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func initConfig(configPath string) (*config.Config, error) {
	err := config.CreateDefaults()
	if err != nil {
		return &config.Config{}, nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return &config.Config{}, nil
	}

	return cfg, err
}

// TODO: Clean this up, doing too many things.
func initChecker(cfg *config.Config) (*pricing.Checker, error) {
	log.Info().Msg("Initializing price feeds...")

	// Register Notifier servers.
	notifying.RegisterNotifier(notifying.Email, email.NewEmailNotifier(cfg.SMTP))

	// Init Coinbase Feed
	feed, err := pricing.NewCoinbaseFeed(cfg.Currency, cfg.Watchlist...)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Coinbase Feed.")

		return &pricing.Checker{}, err
	}

	// TODO: move storage out.
	storage := new(memory.Storage)
	storage.Recipients = []memory.Recipient{
		// TODO: squash commits.
		// {Email: "test@email.com", Type: notifying.Email},
	}

	trigs := []*notifying.Trigger{{ID: 1234, Symbol: "BTC", Threshold: "50_000"}}
	notifier := notifying.NewService(storage)
	checker := pricing.NewChecker(notifier, feed, trigs)

	return checker, nil
}
