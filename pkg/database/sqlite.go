package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func NewSQLiteConnection() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../data.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к SQLite: %w", err)
	}
	db.SetMaxOpenConns(1000)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("(Forum) Не удалось проверить связь с SQLite: %w", err)
	}

	return db, nil
}
