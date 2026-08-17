package migrations

import "embed"

// FS stores all migration files to keep deployment self-contained.
//
//go:embed *.up.sql
var FS embed.FS
