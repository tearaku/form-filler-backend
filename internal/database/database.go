package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"teacup1592/form-filler/internal/schoolForm"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
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

func (db *DB) GetEventInfo(ctx context.Context, id int) (*schoolForm.EventInfo, error) {
	var event EventInfo
	const sql = `SELECT * FROM "Event" WHERE id = $1`
	rows, err := db.DbPool.Query(ctx, sql, id)
	if err == nil {
		defer rows.Close()
		if err := pgxscan.ScanOne(&event, rows); err != nil {
			return nil, err
		}
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
	var attendants []Attendance
	const sql2 = `SELECT * FROM "Attendance" WHERE "eventId" = $1 ORDER BY "userId" ASC`
	rows, err = db.DbPool.Query(ctx, sql2, id)
	if err == nil {
		defer rows.Close()
		if err := pgxscan.ScanAll(&attendants, rows); err != nil {
			return nil, err
		}
		event.Attendants = attendants
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
	return event.dto(), nil
}

func (db *DB) GetProfiles(ctx context.Context, idList []int32) ([]*schoolForm.UserProfile, error) {
	var profileList []*UserProfile
	const sql = `SELECT * FROM "Profile" WHERE "userId" = ($1::int[]) ORDER BY "userId" ASC`
	rows, err := db.DbPool.Query(ctx, sql, idList)
	if err == nil {
		defer rows.Close()
		if err := pgxscan.ScanAll(profileList, rows); err != nil {
			return nil, err
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get event from database: failed, %v\n", err)
		return nil, errors.New("get event from database: failed")
	}
	var faList []*schoolForm.UserProfile
	for _, profile := range profileList {
		p, err := profile.dto()
		if err != nil {
			return nil, err
		}
		faList = append(faList, p)
	}
	return faList, nil
}

func (db *DB) GetMinProfiles(ctx context.Context, idList []int32) ([]*schoolForm.MinProfile, error) {
	var profileList []*MinProfile
	const sql = `SELECT * FROM "MinimalProfile" WHERE "userId" = ($1::int[]) ORDER BY "userId" ASC`
	rows, err := db.DbPool.Query(ctx, sql, idList)
	if err == nil {
		defer rows.Close()
		if err := pgxscan.ScanAll(profileList, rows); err != nil {
			return nil, err
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get event from database: failed, %v\n", err)
		return nil, errors.New("get event from database: failed")
	}
	var faList []*schoolForm.MinProfile
	for _, profile := range profileList {
		faList = append(faList, profile.dto())
	}
	return faList, nil
}

func (db *DB) GetMemberByDept(ctx context.Context, des string) (*schoolForm.MinProfile, error) {
	var member *MinProfile
	const sql = `SELECT * FROM "Department" WHERE "description" LIKE '$1'`
	rows, err := db.DbPool.Query(ctx, sql, des)
	if err == nil {
		defer rows.Close()
		if err := pgxscan.ScanOne(member, rows); err != nil {
			return nil, err
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get member by department from database: failed, %v\n", err)
		return nil, errors.New("get member by department from database: failed")
	}
	return member.dto(), nil
}
