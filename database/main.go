package database

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/Improwised/golang-api/config"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // import mysql if it is used
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbURL string
var err error
var goquReturn *goqu.Database

const (
	postgres = "postgres"
	mysql    = "mysql"
	sqlite3  = "sqlite3"
)

// Connect with database
func Connect(cfg config.DBConfig) (*goqu.Database, error) {
	switch cfg.Dialect {
	case postgres:
		goquReturn, err = postgresDBConnection(cfg)
	case mysql:
		goquReturn, err = mysqlDBConnection(cfg)
	case sqlite3:
		goquReturn, err = sqlite3DBConnection(cfg)
	default:
		panic("No suitable dialect found")
	}
	if err != nil {
		panic(err)
	}
	return goquReturn, nil
}

func sqlite3DBConnection(cfg config.DBConfig) (*goqu.Database, error) {

	if _, err = os.Stat(cfg.SqlliteFileName); err != nil {
		file, err := os.Create(cfg.SqlliteFileName)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
	db, err = sql.Open(sqlite3, "./"+cfg.SqlliteFileName)
	if err != nil {
		return nil, err
	}
	return goqu.New(sqlite3, db), err
}

func mysqlDBConnection(cfg config.DBConfig) (*goqu.Database, error) {
	dbURL = cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.Db
	if db == nil {
		db, err = sql.Open(mysql, dbURL)
		if err != nil {
			return nil, err
		}
		return goqu.New(mysql, db), err
	}
	return goqu.New(mysql, db), err
}

func postgresDBConnection(cfg config.DBConfig) (*goqu.Database, error) {
	dbURL = "postgres://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + "/" + cfg.Db + "?" + cfg.QueryString
	if db == nil {
		db, err = sql.Open(postgres, dbURL)
		if err != nil {
			return nil, err
		}
		return goqu.New(postgres, db), err
	}
	return goqu.New(postgres, db), err
}
