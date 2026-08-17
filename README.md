# wgsh

WireGuard Shell (`wgsh`) is a Go project for managing a WireGuard-based VPN environment with a small and controlled administrative layer.

The repository is still in an early stage. The current codebase focuses on the database layer, schema migrations, and the initial CLI scaffold. It is not yet a complete WireGuard controller or a finished production deployment stack.

## Current scope

- SQLite persistence for interfaces and peers
- versioned migrations with SQLC-generated accessors
- initial CLI entry point
- deployment script scaffolding
- repository structure ready for future VPN and firewall modules

## Project structure

```text
wgssh/
├── cmd/
│   └── wgsh/
│       └── cli.go
├── internal/
│   └── db/
│       ├── basic.go
│       ├── db.go
│       ├── migrate.go
│       ├── migrate_test.go
│       ├── migrations/
│       │   ├── 001_initial.up.sql
│       │   └── embed.go
│       ├── queries/
│       │   └── core.sql
│       └── sqlc/
│           ├── core.sql.go
│           ├── db.go
│           ├── models.go
│           └── querier.go
├── scripts/
│   └── deploy-freebsd.sh
├── templates/
├── go.mod
├── go.sum
├── sqlc.yaml
├── LICENSE
└── README.md
```

## Data model

The project currently defines a simple database model centered on:

- `interfaces`
  - `id`
  - `name`
  - `private_key`
  - `public_key`
  - `listen_port`
  - `ip_address`
  - timestamps

- `peers`
  - `id`
  - `interface_id`
  - `endpoint`
  - `name`
  - `public_key`
  - `private_key`
  - `allowed_ips`
  - `allow_intercommunication`
  - `latest_handshake`
  - timestamps

The schema is handled through migrations and SQLC-based generated code.

## Requirements

- Go 1.25+ (module currently points to Go 1.26.5)
- FreeBSD 15+ as the main target environment
- WireGuard and PF for runtime use in the future
- SQLite

## Build and test

```bash
go build ./cmd/wgsh
go test ./...
```

## Status

This repository is intentionally scoped to the foundation layer. The next meaningful milestones are:

- implement the WireGuard runtime module
- implement the firewall layer for FreeBSD/Linux backends
- define the operational CLI flow and SSH-controlled execution model
- generate client configuration artifacts and peer management logic

## License

This project is licensed under the BSD 3-Clause License.

Copyright (c) 2026, Psykka
All rights reserved.

See [LICENSE](LICENSE) for the full text.
