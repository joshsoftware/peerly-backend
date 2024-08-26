package main

// @APITitle Main
// @APIDescription Main API for Microservices in Go!

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"os"
	"strconv"

	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/api"
	"github.com/joshsoftware/peerly-backend/internal/app"
	"github.com/joshsoftware/peerly-backend/internal/app/cronjob"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	script "github.com/joshsoftware/peerly-backend/scripts"
	"github.com/rs/cors"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	config.Load()

	cliApp := cli.NewApp()
	cliApp.Name = config.AppName()
	cliApp.Version = "1.0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start server",
			Action: func(c *cli.Context) error {
				return startApp()
			},
		},
		{
			Name:  "create_migration",
			Usage: "create migration file",
			Action: func(c *cli.Context) error {
				return repository.CreateMigrationFile(c.Args().Get(0))
			},
		},
		{
			Name:  "migrate",
			Usage: "run db migrations",
			Action: func(c *cli.Context) error {
				return repository.RunMigrations()
			},
		},
		{
			Name:      "rollback",
			Usage:     "rollback migrations [step (int)]",
			ArgsUsage: "[step (int)]",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return errors.New("migration step is required")
				}
				return repository.RollbackMigrations(c.Args().Get(0))
			},
		},
		{
			Name:  "seed",
			Usage: "seed data in database",
			Action: func(c *cli.Context) error {
				return repository.SeedData()
			},
		},
		{
			Name:  "loadUsers",
			Usage: "load peerly users from intranet",
			Action: func(c *cli.Context) error {
				return script.LoadUserScript()
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func startApp() (err error) {

	// Context for main function
	ctx := context.Background()
	lg,err := log.SetupLogger()
	if err != nil{
		logger.Error("logger setup failed ",err.Error())
		return err
	}
	log.Info(ctx,"Starting Peerly Application...")
	defer log.Info(ctx,"Shutting Down Peerly Application...")
	//initialize database
	dbInstance, err := repository.InitializeDatabase()
	if err != nil {
		log.Error(ctx,"Database init failed")
		return err
	}

	//cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
	})

	//initialize service dependencies
	services := app.NewService(dbInstance)

	// Initializing Cron Job
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Error(ctx, "scheduler creation failed with error: %s", err.Error())
		return err
	}

	cronjob.InitializeJobs(services.AppreciationService, services.UserService, scheduler)
	// defer func() {
    //     if err := scheduler.Shutdown(); err != nil {
    //         log.Error(ctx, "Scheduler shutdown failed: %s", err.Error())
    //     }
    // }()
	defer scheduler.Shutdown()
	//initialize router
	router := api.NewRouter(services)

	// Negroni logger setup
	negroniLogger := negroni.NewLogger()
	negroniLogger.ALogger = lg

	// Initialize web server
	server := negroni.New(negroniLogger)
	server.Use(c)
	server.UseHandler(router)

	port := config.AppPort()
	addr := fmt.Sprintf(":%s", strconv.Itoa(port))
	server.Run(addr)
	return
}
