package datasrc

import (
	"context"
	"errors"
	"log"
	"os"
	"teacup1592/form-filler/internal/schoolForm"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	Postgres *pgx.Conn
}

func (db *DB) GetEventInfo(ctx context.Context, params schoolForm.GetEventInfoParams) (*schoolForm.EventInfo, error) {
	var event schoolForm.EventInfo
	const sql = `SELECT * FROM "Event" WHERE id = $1`
	rows, err := db.Postgres.Query(ctx, sql, params.EventID)
	if err == nil {
		defer rows.Close()
		pgxscan.ScanOne(&event, rows)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		log.Printf("Get event from database: failed, %v\n", err)
		return nil, errors.New("get event from database: failed")
	}
	return &event, nil
}

// Creates a single connection
// TODO: if needed, switch to connection pool
func ConnectToDb() *DB {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	log.Println("Connection to database: ok!")
	// defer conn.Close(context.Background())
	return &DB{Postgres: conn}
}
