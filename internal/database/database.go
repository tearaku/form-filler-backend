package database

import (
	"context"
	"fmt"
	"teacup1592/form-filler/internal/schoolForm"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	DbPool *pgxpool.Pool
}

func NewDBPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgxpool connection error: %w", err)
	}
	return pool, nil
}

func (db *DB) GetEventInfo(ctx context.Context, params schoolForm.GetEventInfoParams) (*schoolForm.EventInfo, error) {
	return nil, nil
}
