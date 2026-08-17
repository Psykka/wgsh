# wgsh

WireGuard Shell (wgsh) is a minimalist and autonomous orchestrator for provisioning and managing WireGuard VPN tunnels in a FreeBSD infrastructure, with a focus on security, logical isolation, and restricted SSH execution.

The project is designed to operate as an administrative gateway, where the operator connects over SSH to a dedicated user and the Go binary acts as the local controller for VPN state, the database, and network rules.

## Objective

- Automate the creation and maintenance of WireGuard interfaces;
- Manage peers, addresses, public/private keys, and communication permissions;
- Rewrite PF firewall rules atomically when the network state changes;
- Apply default isolation so that only explicitly allowed traffic is permitted between clients;
- Provide a simple CLI for secure remote administration.

## Architecture Overview

The proposed operating flow is as follows:

1. The administrator accesses the FreeBSD host over SSH using key-based authentication and a dedicated user.
2. The service user (`vpn@freebsd-vpn`) does not have an interactive shell and is redirected directly to the `wgsh` binary through `ForceCommand` in `sshd_config`.
3. The program runs with minimal privileges for the regular user and only elevates the strictly necessary OS commands through `doas`.
4. The application uses a local SQLite database to persist interfaces, peers, keys, and permissions.
5. The PF firewall is updated in real time to reflect the isolation and permitted communication policy.

## Technology Stack

- Go
- WireGuard
- FreeBSD 15
- PF (Packet Filter)
- Linux + `iptables`/`nftables` (future compatibility path)
- SQLite
- `modernc.org/sqlite` (pure Go driver)
- `doas` for controlled privilege escalation

## Security Principles

- A dedicated user for SSH access;
- Locked shell for service users;
- Binary execution through `ForceCommand` instead of an interactive shell;
- Surgical escalation only for sensitive actions;
- Atomic rewrite of firewall rules;
- Default network isolation (`block in on wg0 from vpn_net to vpn_net`), with exceptions controlled per peer;
- Explicit control of intercommunication between peers via a boolean permission flag.

## Planned / Core Features

- Registration and management of WireGuard interfaces;
- Generation of Curve25519 keys in memory;
- Local persistence in SQLite;
- Registration of peers with endpoint, allowed IPs, and intercommunication options;
- Traffic, handshake, and connectivity monitoring;
- Client configuration generation for distribution via QR code and plain text in the terminal;
- Dynamic updates of PF rules and synchronization with `pfctl`;
- Future compatibility with Linux, including isolation rules via `iptables`;
- Integration with `wg syncconf` to synchronize the kernel interface.

## Repository Structure

```text
wgssh/
├── cmd/
│   └── wgsh/
│       └── cli.go             # Main administrator CLI
├── internal/
│   └── db/
│       ├── basic.go           # Data-layer helper accessors
│       ├── db.go              # SQLite initialization and Store setup
│       ├── migrate.go         # Schema migration logic
│       ├── migrate_test.go    # Migration tests
│       ├── migrations/
│       │   ├── 001_initial.up.sql
│       │   └── embed.go
│       ├── queries/
│       │   └── core.sql       # Core SQL queries
│       └── sqlc/
│           ├── core.sql.go
│           ├── db.go
│           ├── models.go
│           └── querier.go
├── scripts/
│   └── deploy-freebsd.sh     # Deployment script for FreeBSD
├── templates/                # Templates for configuration and artifact generation
├── go.mod                    # Go module definition
├── go.sum
├── sqlc.yaml                 # sqlc configuration
└── README.md
```

## Data Model

The current SQLite schema includes the main entities:

- `interfaces`
  - `id`
  - `name`
  - `private_key`
  - `public_key`
  - `listen_port`
  - `ip_address`
  - creation and update timestamps

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
  - creation and update timestamps

The database is organized with versioned migrations, using the `schema_migrations` table to ensure consistency and traceability of schema changes.

## Execution Flow

The intended operational flow is:

1. Secure SSH connection using a dedicated user and public key authentication.
2. Login redirected to the `wgsh` application via `ForceCommand`.
3. Loading of the local environment and SQLite database.
4. Reading and updating interfaces and peers.
5. Client configuration generation.
6. Updating `pf.conf` and reloading the firewall.
7. Synchronization with the WireGuard kernel using `wg syncconf`.
8. Display of the ASCII QR code and configuration text to the client.

## Requirements

### Runtime environment

- FreeBSD 15+ (primary target)
- Linux with WireGuard + `iptables`/`nftables` (future plan)
- WireGuard enabled in the kernel
- PF enabled on FreeBSD
- `doas` configured
- SSH access with key authentication and shell restrictions

### Development

- Go 1.25+ (the current module points to Go 1.26.5)
- Go module-based project structure
- `sqlc` for SQL code generation

## How to Build

```bash
go build ./cmd/wgsh
```

## How to Test

```bash
go test ./...
```

## Security and Operations

> This project was designed with a zero-trust model and source/destination isolation. In practice, peers on a VPN network should not communicate with each other by default; any explicit communication must be authorized and reflected in the database and firewall rules.

> Future plan: expand project compatibility for Linux by using `iptables` to apply policies equivalent to PF on FreeBSD. The business layer and state persistence will follow the same access-control rules, but the firewall implementation will be adapted to the target operating system backend.

## Current Project Status

The repository already contains the core building blocks for the SQLite data model and migration system, along with the initial CLI structure. The project is still evolving toward the WireGuard management layer, PF integration, and SSH-controlled authentication and execution flows.

## License

This project is licensed under the BSD 3-Clause License.

Copyright (c) 2026, Psykka
All rights reserved.

See [LICENSE](LICENSE) for the full text.

## Recommended Next Steps

- Define the `main.go` layout and application entry point;
- Implement the WireGuard module (`internal/vpn`);
- Implement the firewall module (`internal/firewall`), with backend-specific support for PF on FreeBSD and `iptables` on Linux;
- Finalize the interactive administration CLI;
- Validate the full provisioning flow in a real FreeBSD environment;
- Validate Linux compatibility with `iptables` and peer-isolation rules;
- Document the deployment process and the configuration of `sshd_config` and `doas`.
