package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	logger "github.com/sirupsen/logrus"

	// Import PostgreSQL database driver
	_ "github.com/lib/pq"

	// For database migrations
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	// golang-migrate reads migrations from sources and applies them in correct order to a database.
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	dbDriver = "postgres"
)

var (
	Sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

// InitializeDatabase initialize database and return database instance
func InitializeDatabase() (db *sqlx.DB, err error) {
	uri := config.ReadEnvString(constants.DBURI)

	conn, err := sqlx.Connect(dbDriver, uri)
	if err != nil {
		logger.Errorf("Cannot initialize database: %v",err)
		return
	}
	return conn, nil
}

// RunMigrations - runs all database migrations (see ../migrtions/*.up.sql)
func RunMigrations() (err error) {
	uri := config.ReadEnvString(constants.DBURI)

	db, _ := sql.Open(dbDriver, uri)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Errorf("failure to create driver obj: %v",err)
		return apperrors.FailedToCreateDriver
	}

	m, err := migrate.NewWithDatabaseInstance(getMigrationPath(), dbDriver, driver)
	if err != nil {
		logger.Errorf("migrate failure: %v",err)
		return apperrors.MigrationFailure
	}

	err = m.Up()
	if err == migrate.ErrNoChange || err == nil {
		err = nil
		return
	}

	return
}

// CreateMigrationFile - Creates a boilerplate *.sql files for a database migration
func CreateMigrationFile(filename string) (err error) {
	if len(filename) == 0 {
		err = errors.New("filename is not provided")
		return
	}
	
	timeStamp := time.Now().Unix()
	upMigrationFilePath := fmt.Sprintf("%s/%d_%s.up.sql", config.ReadEnvString(constants.MigrationFolderPath), timeStamp, filename)
	downMigrationFilePath := fmt.Sprintf("%s/%d_%s.down.sql", config.ReadEnvString(constants.MigrationFolderPath), timeStamp, filename)
	err = createFile(upMigrationFilePath)
	if err != nil {
		return
	}

	err = createFile(downMigrationFilePath)
	if err != nil {
		os.Remove(upMigrationFilePath)
		return
	}

	logger.WithFields(logger.Fields{
		"up":   upMigrationFilePath,
		"down": downMigrationFilePath,
	}).Info("Created migration files")

	return
}

// RollbackMigrations - Used to run the "down" database migrations in ../migrations/*.down.sql
func RollbackMigrations(s string) (err error) {
	uri := config.ReadEnvString(constants.DBURI)

	steps, err := strconv.Atoi(s)
	if err != nil {
		return
	}

	m, err := migrate.New(getMigrationPath(), uri)
	if err != nil {
		return
	}

	err = m.Steps(-1 * steps)
	if err == migrate.ErrNoChange || err == nil {
		err = nil
		return
	}

	return
}

func createFile(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	err = f.Close()
	return
}

func getMigrationPath() string {
	return fmt.Sprintf("file://%s", config.ReadEnvString(constants.MigrationFolderPath))
}
