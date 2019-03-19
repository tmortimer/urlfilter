package connectors

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tmortimer/urlfilter/config"
)

const CREATE_DB string = "CREATE DATABASE IF NOT EXISTS URLFilter"

const USE_DB string = "USE URLFilter"

// High Performance MYSQL from O'REilly suggested the CRC as an index approach.
const CREATE_URL_TABLE = "CREATE TABLE IF NOT EXISTS crcurls (" +
	"id int unsigned NOT NULL auto_increment," +
	"url varchar(2050) NOT NULL," +
	"url_crc int unsigned NOT NULL DEFAULT 0," +
	"PRIMARY KEY(id)," +
	"INDEX(url_crc)" +
	")"

// Holds the MySQL connection pool and executes commands against MySQL.
type MySQL struct {
	// MySQL connection pool.
	db *sql.DB

	// MySQL specific config.
	config config.MySQL
}

// Create a new MySQL connector and setup the MySQL connection pool.
func NewMySQL(config config.MySQL) (*MySQL, error) {
	connector := &MySQL{
		config: config,
	}

	err := connector.ConfigureMySQL()
	if err != nil {
		r.db.Close() // Probably don't need this.
		return nil, err
	}
	return connector, nil
}

// Setup MySQL schema.
func (r *MySQL) ConfigureMySQL() error {
	dsn := r.config.Username + ":" + r.config.Password +
		"@tcp(" + r.config.Host + ":" + r.config.Port + ")/"

	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		return err
	}

	_, err = db.Exec(CREATE_DB)
	if err != nil {
		return err
	}

	_, err = db.Exec(USE_DB)
	if err != nil {
		return err
	}

	_, err = db.Exec(CREATE_URL_TABLE)
	if err != nil {
		return err
	}

	// The USE_DB statement above won't apply to all connections
	// in the existing pool, so we actually want to open a new one
	// and close the old one.
	r.db, err = sql.Open("mysql", dsn+"URLFilter")

	return err
}

// Check if the URL is in MySQL.
func (r *MySQL) ContainsURL(url string) (bool, error) {

	return false, nil
}

// Add the URL to the MySQL. Only used if this DB is being used as a cache.
func (r *MySQL) AddURL(url string) error {

	return nil
}

// Return the name MySQL for logging.
func (r *MySQL) Name() string {
	return "MySQL"
}
