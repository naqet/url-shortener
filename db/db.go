package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/naqet/url-shortener/models"
)

type DB struct {
	connection *sql.DB
}

func NewDB(path string) (*DB, error) {
	connection, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

    db := DB{connection};

    err = db.Init();

    if err != nil {
        return nil, err;
    }

	return &db, nil
}

func (db *DB) Close() error {
	if db.connection == nil {
		return errors.New("Connection to the DB is not available")
	}

	db.connection.Close()
	return nil
}

func (db *DB) Init() error {
    _, err := db.connection.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            key TEXT NOT NULL UNIQUE,
            original_url TEXT NOT NULL
        );
    `)

	return err
}

func (db *DB) SaveURL(key, originalUrl string) error {
    _, err := db.connection.Exec("INSERT INTO urls (key, original_url) VALUES (?, ?)", key, originalUrl);
    return err;
}

func (db *DB) GetURL(key string) (models.ShortUrl, error) {
    url := models.ShortUrl{}
    err := db.connection.QueryRow("SELECT * FROM urls WHERE key = ?", key).Scan(&url.Id, &url.Key, &url.OriginalUrl);
    return url, err;
}
