# wgsh

WireGuard Shell (`wgsh`) is a minimalist and autonomous orchestrator for provisioning and managing WireGuard VPN tunnels in a FreeBSD infrastructure.

This repository is still in an early stage. It currently contains the foundation for the data layer, migrations, and initial CLI scaffolding, but it is not yet a complete production VPN orchestrator.

## Objective

Provide a minimal and controlled way to manage WireGuard-related state and administrative operations from a local Go application, with a future path toward FreeBSD-oriented firewall and SSH-controlled execution.

## Stack

- Go
- SQLite
- SQLC
- WireGuard
- FreeBSD / PF (planned runtime target)

## Current status

The repository currently contains:

- initial database schema and migration logic
- generated SQLC access layer
- SQLite helpers and basic storage pattern
- a CLI shell entry point

What is still missing or intentionally not implemented yet:

- WireGuard runtime management
- firewall automation
- SSH access control flow
- peer/client configuration generation
- production deployment logic

## Build and test

```bash
go build ./cmd/wgsh
go test ./...
```

## License

This project is licensed under the BSD 3-Clause License.

Copyright (c) 2026, Psykka
All rights reserved.

See [LICENSE](LICENSE) for the full text.
