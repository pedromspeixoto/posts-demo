package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pedromspeixoto/posts-api/internal/config"
	"github.com/pedromspeixoto/posts-api/internal/pkg/filepath"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbDeps struct {
	fx.In

	Config *config.Config
}

type DbClient struct {
	dbDeps
}

func NewDbClient(deps dbDeps) (*gorm.DB, error) {

	dsn := deps.Config.MySQLUrl()

	// create db and run migrations on startup if not production
	if deps.Config.Environment != config.EnvironmentProduction {
		// connect to default instance (to create db)
		dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			deps.Config.MySQLUser,
			deps.Config.MySQLPassword,
			deps.Config.MySQLHost,
			deps.Config.MySQLPort,
			"mysql",
		)

		err := createDb(dbUrl, deps.Config.MySQLDBName)
		if err != nil {
			log.Fatalf("error creating database: %v", err)
		}

		err = migrateDb(dsn)
		if err != nil {
			log.Fatalf("error migrating database: %v", err)
		}
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err = sqldb.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func createDb(url, name string) error {
	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", name)); err != nil {
		return errors.Wrap(err, "failed to create database")
	}
	return nil
}

func migrateDb(url string) error {
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println("error opening database")
		return err
	}

	goose.SetBaseFS(nil)
	goose.SetTableName("goose_db_version")

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	if err := goose.Up(db,
		fmt.Sprintf("%s/migrations", filepath.ProjectRootDir()),
		goose.WithAllowMissing()); err != nil {
		return err
	}

	return nil
}
