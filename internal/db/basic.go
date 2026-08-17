package db

import (
	"context"
	"database/sql"

	"wgsh/internal/db/sqlc"
)

type CreateInterfaceInput struct {
	Name       string
	PrivateKey string
	PublicKey  string
	ListenPort int64
	IPAddress  string
}

type UpdateInterfaceInput struct {
	ID         int64
	PrivateKey string
	PublicKey  string
	ListenPort int64
	IPAddress  string
}

type CreatePeerInput struct {
	InterfaceID             int64
	Endpoint                string
	Name                    string
	PublicKey               string
	PrivateKey              string
	AllowedIPs              string
	AllowIntercommunication bool
}

type UpdatePeerInput struct {
	ID                      int64
	Endpoint                string
	Name                    string
	PublicKey               string
	PrivateKey              string
	AllowedIPs              string
	AllowIntercommunication bool
}

func (s *Store) CreateInterface(ctx context.Context, in CreateInterfaceInput) (sqlc.Interface, error) {
	return s.Queries.CreateInterface(ctx, sqlc.CreateInterfaceParams{
		Name:       in.Name,
		PrivateKey: in.PrivateKey,
		PublicKey:  in.PublicKey,
		ListenPort: in.ListenPort,
		IpAddress:  in.IPAddress,
	})
}

func (s *Store) GetInterfaceByID(ctx context.Context, id int64) (sqlc.Interface, error) {
	return s.Queries.GetInterfaceByID(ctx, id)
}

func (s *Store) GetInterfaceByName(ctx context.Context, name string) (sqlc.Interface, error) {
	return s.Queries.GetInterfaceByName(ctx, name)
}

func (s *Store) ListInterfaces(ctx context.Context) ([]sqlc.Interface, error) {
	return s.Queries.ListInterfaces(ctx)
}

func (s *Store) UpdateInterface(ctx context.Context, in UpdateInterfaceInput) error {
	return s.Queries.UpdateInterface(ctx, sqlc.UpdateInterfaceParams{
		ID:         in.ID,
		PrivateKey: in.PrivateKey,
		PublicKey:  in.PublicKey,
		ListenPort: in.ListenPort,
		IpAddress:  in.IPAddress,
	})
}

func (s *Store) DeleteInterfaceByID(ctx context.Context, id int64) error {
	return s.Queries.DeleteInterfaceByID(ctx, id)
}

func (s *Store) CreatePeer(ctx context.Context, in CreatePeerInput) (sqlc.Peer, error) {
	return s.Queries.CreatePeer(ctx, sqlc.CreatePeerParams{
		InterfaceID:             in.InterfaceID,
		Endpoint:                in.Endpoint,
		Name:                    in.Name,
		PublicKey:               in.PublicKey,
		PrivateKey:              in.PrivateKey,
		AllowedIps:              in.AllowedIPs,
		AllowIntercommunication: boolToNullInt64(in.AllowIntercommunication),
	})
}

func (s *Store) GetPeerByID(ctx context.Context, id int64) (sqlc.Peer, error) {
	return s.Queries.GetPeerByID(ctx, id)
}

func (s *Store) ListPeersByInterfaceID(ctx context.Context, interfaceID int64) ([]sqlc.Peer, error) {
	return s.Queries.ListPeersByInterfaceID(ctx, interfaceID)
}

func (s *Store) UpdatePeer(ctx context.Context, in UpdatePeerInput) error {
	return s.Queries.UpdatePeer(ctx, sqlc.UpdatePeerParams{
		ID:                      in.ID,
		Endpoint:                in.Endpoint,
		Name:                    in.Name,
		PublicKey:               in.PublicKey,
		PrivateKey:              in.PrivateKey,
		AllowedIps:              in.AllowedIPs,
		AllowIntercommunication: boolToNullInt64(in.AllowIntercommunication),
	})
}

func (s *Store) SetPeerLatestHandshake(ctx context.Context, id int64, latestHandshake int64) error {
	return s.Queries.SetPeerLatestHandshake(ctx, sqlc.SetPeerLatestHandshakeParams{
		ID:              id,
		LatestHandshake: latestHandshake,
	})
}

func (s *Store) DeletePeerByID(ctx context.Context, id int64) error {
	return s.Queries.DeletePeerByID(ctx, id)
}

func boolToNullInt64(v bool) sql.NullInt64 {
	if v {
		return sql.NullInt64{Int64: 1, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: true}
}
