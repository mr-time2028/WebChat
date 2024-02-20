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

// pingDB ping to the database to ensure database is open
func pingDB(d *sql.DB) error {
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

	err = pingDB(sdb)
	if err != nil {
		return nil, err
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return nil, err
	}

	return &DB{
		GormDB: db,
		SqlDB:  sdb,
	}, nil
}

func ConnectTestSQL() (*DB, error) {
	testDBName := helpers.GetEnvOrDefaultString("TEST_POSTGRES_NAME", "")
	testDBUser := helpers.GetEnvOrDefaultString("TEST_POSTGRES_USER", "")
	testDBPass := helpers.GetEnvOrDefaultString("TEST_POSTGRES_PASSWORD", "")
	testDBHost := helpers.GetEnvOrDefaultString("TEST_POSTGRES_HOST", "localhost")
	testDBPort := helpers.GetEnvOrDefaultString("TEST_POSTGRES_PORT", "5432")
	testDBSSL := helpers.GetEnvOrDefaultString("TEST_POSTGRES_SSL", "disable")
	testDBZone := helpers.GetEnvOrDefaultString("TEST_POSTGRES_ZONE", "Asia/Tehran")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		testDBHost, testDBUser, testDBPass, testDBName, testDBPort, testDBSSL, testDBZone)

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	sdb, _ := db.DB()

	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return nil, err
	}

	return &DB{
		GormDB: db,
		SqlDB:  sdb,
	}, nil
}

// GetAllTables gather all tables name that exists in the database
func (db *DB) GetAllTables() ([]string, error) {
	tables, err := db.GormDB.Migrator().GetTables()
	if err != nil {
		return nil, err
	}

	return tables, nil
}

// DropAllTables drop all tables in the database
func (db *DB) DropAllTables() error {
	tables, err := db.GetAllTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if err = db.GormDB.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
