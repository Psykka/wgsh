package db

import (
	"database/sql"
	"fmt"

	"wgsh/internal/db/sqlc"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB      *sql.DB
	Queries *sqlc.Queries
}

func Open(path string) (*Store, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, stmt := range pragmas {
		if _, execErr := conn.Exec(stmt); execErr != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("apply pragma %q: %w", stmt, execErr)
		}
	}

	if err := Migrate(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &Store{
		DB:      conn,
		Queries: sqlc.New(conn),
	}, nil
}

func (s *Store) Close() error {
	if s == nil || s.DB == nil {
		return nil
	}
	return s.DB.Close()
}
