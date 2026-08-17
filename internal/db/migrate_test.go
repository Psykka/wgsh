package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrateAppliesSchema(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "wgsh.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	if err := Migrate(conn); err != nil {
		t.Fatalf("first migrate failed: %v", err)
	}

	if err := Migrate(conn); err != nil {
		t.Fatalf("second migrate failed: %v", err)
	}

	for _, table := range []string{"interfaces", "peers", "schema_migrations"} {
		var found string
		err := conn.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ? LIMIT 1`, table).Scan(&found)
		if err != nil {
			t.Fatalf("table %s not found: %v", table, err)
		}
	}
}
