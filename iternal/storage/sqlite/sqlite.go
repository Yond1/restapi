package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("alias %s already exists", alias)
		}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	var url string
	err := s.db.QueryRow("SELECT url FROM url WHERE alias = ?", alias).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("alias %s not found", alias)
		}
		return "", fmt.Errorf("failed to get url: %w", err)
	}
	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	_, err := s.db.Exec("DELETE FROM url WHERE alias = ?", alias)
	if err != nil {
		return fmt.Errorf("failed to delete url: %w", err)
	}
	return nil
}
