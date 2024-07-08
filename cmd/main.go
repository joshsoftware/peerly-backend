package main

// @APITitle Main
// @APIDescription Main API for Microservices in Go!

import (
	"errors"
	"fmt"
	"net/http"

	"os"
	"strconv"

	"github.com/joshsoftware/peerly-backend/internal/api"
	"github.com/joshsoftware/peerly-backend/internal/app"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/repository"
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
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func startApp() (err error) {

	logger.Info("Starting Peerly Application...")
	defer logger.Info("Shutting Down Peerly Application...")
	//initialize database
	dbInstance, err := repository.InitializeDatabase()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Database init failed")
		return
	}

	//cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
	})

	//initialize service dependencies
	services := app.NewService(dbInstance)

	//initialize router
	router := api.NewRouter(services)

	// init web server
	server := negroni.Classic()
	server.Use(c)
	server.UseHandler(router)

	port := config.AppPort()
	addr := fmt.Sprintf(":%s", strconv.Itoa(port))
	server.Run(addr)
	return
}
