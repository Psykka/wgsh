package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"wgsh/internal/db/migrations"
)

func Migrate(conn *sql.DB) error {
	if conn == nil {
		return fmt.Errorf("nil sqlite connection")
	}

	if err := ensureMigrationsTable(conn); err != nil {
		return err
	}

	migrationFiles, err := migrationNames()
	if err != nil {
		return err
	}

	for _, name := range migrationFiles {
		if err := applyMigration(conn, name); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(conn *sql.DB) error {
	_, err := conn.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            checksum TEXT NOT NULL,
            applied_at INTEGER NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	return nil
}

func migrationNames() ([]string, error) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return names, nil
}

func applyMigration(conn *sql.DB, name string) error {
	migrationSQL, err := fs.ReadFile(migrations.FS, name)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", name, err)
	}

	hash := sha256.Sum256(migrationSQL)
	checksum := hex.EncodeToString(hash[:])

	var existingChecksum string
	err = conn.QueryRow(`SELECT checksum FROM schema_migrations WHERE version = ? LIMIT 1`, name).Scan(&existingChecksum)
	switch {
	case err == nil:
		if existingChecksum != checksum {
			return fmt.Errorf("migration %s checksum mismatch", name)
		}
		return nil
	case err != sql.ErrNoRows:
		return fmt.Errorf("query migration %s: %w", name, err)
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", name, err)
	}

	if _, err := tx.Exec(string(migrationSQL)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("apply migration %s: %w", name, err)
	}

	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version, checksum, applied_at) VALUES (?, ?, strftime('%s','now'));`,
		name,
		checksum,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %s: %w", name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", name, err)
	}

	return nil
}
