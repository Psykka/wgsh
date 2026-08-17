-- name: CreateInterface :one
INSERT INTO interfaces (
	name,
	private_key,
	public_key,
	listen_port,
	ip_address
)
VALUES (?, ?, ?, ?, ?)
RETURNING id, name, private_key, public_key, listen_port, ip_address, created_at, updated_at;

-- name: GetInterfaceByID :one
SELECT id, name, private_key, public_key, listen_port, ip_address, created_at, updated_at
FROM interfaces
WHERE id = ?
LIMIT 1;

-- name: GetInterfaceByName :one
SELECT id, name, private_key, public_key, listen_port, ip_address, created_at, updated_at
FROM interfaces
WHERE name = ?
LIMIT 1;

-- name: ListInterfaces :many
SELECT id, name, private_key, public_key, listen_port, ip_address, created_at, updated_at
FROM interfaces
ORDER BY name ASC;

-- name: UpdateInterface :exec
UPDATE interfaces
SET
	private_key = ?,
	public_key = ?,
	listen_port = ?,
	ip_address = ?,
	updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteInterfaceByID :exec
DELETE FROM interfaces
WHERE id = ?;

-- name: CreatePeer :one
INSERT INTO peers (
	interface_id,
	endpoint,
	name,
	public_key,
	private_key,
	allowed_ips,
	allow_intercommunication
)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id, interface_id, endpoint, name, public_key, private_key, allowed_ips, allow_intercommunication, latest_handshake, created_at, updated_at;

-- name: GetPeerByID :one
SELECT id, interface_id, endpoint, name, public_key, private_key, allowed_ips, allow_intercommunication, latest_handshake, created_at, updated_at
FROM peers
WHERE id = ?
LIMIT 1;

-- name: ListPeersByInterfaceID :many
SELECT id, interface_id, endpoint, name, public_key, private_key, allowed_ips, allow_intercommunication, latest_handshake, created_at, updated_at
FROM peers
WHERE interface_id = ?
ORDER BY name ASC;

-- name: UpdatePeer :exec
UPDATE peers
SET
	endpoint = ?,
	name = ?,
	public_key = ?,
	private_key = ?,
	allowed_ips = ?,
	allow_intercommunication = ?,
	updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetPeerLatestHandshake :exec
UPDATE peers
SET
	latest_handshake = ?,
	updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeletePeerByID :exec
DELETE FROM peers
WHERE id = ?;
