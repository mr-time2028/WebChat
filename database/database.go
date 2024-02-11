package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mr-time2028/WebChat/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB contains sql and nosql dbs
type DB struct {
	GormDB *gorm.DB
	SqlDB  *sql.DB
}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

// getDSN return dsn string for connection to the database
func getDSN() string {
	dbName := helpers.GetEnvOrDefaultString("POSTGRES_NAME", "")
	dbUser := helpers.GetEnvOrDefaultString("POSTGRES_USER", "")
	dbPass := helpers.GetEnvOrDefaultString("POSTGRES_PASSWORD", "")
	dbHost := helpers.GetEnvOrDefaultString("POSTGRES_HOST", "localhost")
	dbPort := helpers.GetEnvOrDefaultString("POSTGRES_PORT", "5432")
	dbSSL := helpers.GetEnvOrDefaultString("POSTGRES_SSL", "disable")
	dbZone := helpers.GetEnvOrDefaultString("POSTGRES_ZONE", "Asia/Tehran")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSL, dbZone)

	return dsn
}

// openDB open the database with dsn from getDSN
func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, err
}

// testDB ping to the database to ensure database is open
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// ConnectSQL get dsn and open db and return DB instance
func ConnectSQL() (*DB, error) {
	dsn := getDSN()

	db, err := openDB(dsn)
	if err != nil {
		return nil, err
	}

	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxOpenConns(maxOpenDBConn)
	sdb.SetMaxIdleConns(maxIdleDBConn)
	sdb.SetConnMaxLifetime(maxDBLifetime)

	err = testDB(sdb)
	if err != nil {
		return nil, err
	}

	return &DB{
		GormDB: db,
		SqlDB:  sdb,
	}, nil
}
