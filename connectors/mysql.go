package connectors

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tmortimer/urlfilter/config"
	"hash/crc32"
)

const CREATE_DB string = "CREATE DATABASE IF NOT EXISTS URLFilter"

const USE_DB string = "USE URLFilter"

// High Performance MYSQL from O'REilly suggested the CRC as an index approach.
// This is from Chapter 3: Schema Optimization and Indexing. The idea here is
// that the CRC32 index is a lot faster to lookup than a string based index on the
// url. It's possible to run into collisions, but that still likely only results
// in a few rows to scan through.
//
// With a greater understanding of the sample size, and ideally with some sample
// data one could evaluate if it's worth using CRC64, or even something else.

const CREATE_URL_TABLE = "CREATE TABLE IF NOT EXISTS crcurls (" +
	"id int unsigned NOT NULL auto_increment," +
	"url varchar(2050) NOT NULL," +
	"url_crc int unsigned NOT NULL DEFAULT 0," +
	"PRIMARY KEY(id)," +
	"INDEX(url_crc)" +
	")"

const SELECT_URL = "SELECT EXISTS(SELECT 1 FROM crcurls WHERE url_crc=? AND url=?)"

const ADD_URL = "INSERT INTO crcurls (url_crc, url) VALUES (?, ?)"

const SELECT_RANGE = "SELECT url FROM crcurls WHERE id BETWEEN ? AND ?"

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
	exists := false

	row := r.db.QueryRow(SELECT_URL, crc32.ChecksumIEEE([]byte(url)), url)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// Add the URL to the MySQL. Only used if this DB is being used as a cache.
func (r *MySQL) AddURL(url string) error {
	_, err := r.db.Exec(ADD_URL, crc32.ChecksumIEEE([]byte(url)), url)
	return err
}

// Return the name MySQL for logging.
func (r *MySQL) Name() string {
	return "MySQL"
}

func (r *MySQL) GetURLPage(start int, number int) ([]string, error) {
	urls := make([]string, 0, number)
	rows, err := r.db.Query(SELECT_RANGE, start, start + number - 1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(&urls[i])
		if err != nil {
			return nil, err
		}
		i++
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return urls, nil
}
