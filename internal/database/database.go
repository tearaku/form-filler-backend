package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"teacup1592/form-filler/internal/schoolForm"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgtype"
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

// TODO: I need meaningful error msg to trace things @@
func (db *DB) GetEventInfo(ctx context.Context, id int) (*schoolForm.EventInfo, error) {
	var event EventInfo
	const sql = `SELECT * FROM "Event" WHERE id = $1`
	rows, err := db.DbPool.Query(ctx, sql, id)
	if err == nil {
		if err = pgxscan.ScanOne(&event, rows); err != nil {
			log.Printf("Scanning into db.EventInfo failed: %v\n", err)
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get event from database: failed, %v\n", err)
		return nil, errors.New("failed to get event from database")
	}
	var attendants []Attendance
	const sql2 = `SELECT * FROM "Attendance" WHERE "eventId" = $1 ORDER BY "userId" ASC`
	rows, err = db.DbPool.Query(ctx, sql2, id)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get event's attendance from database: failed, %v\n", err)
		return nil, errors.New("failed to get event info 'attendance' from database")
	}
	defer rows.Close()
	rs := pgxscan.NewRowScanner(rows)
	isEmpty := true
	for rows.Next() {
		var a Attendance
		if err := rs.Scan(&a); err != nil {
			log.Printf("Scanning attendance rows from db failed: %\v\n", err)
			return nil, errors.New("failed to parse event info 'attendance' from database")
		}
		attendants = append(attendants, a)
		isEmpty = false
	}
	if isEmpty {
		log.Printf("Empty Attendance result set\n")
		return nil, errors.New("no event attendance were fetched")
	}
	event.Attendants = attendants
	return event.dto(), nil
}

func (db *DB) GetProfiles(ctx context.Context, idList []int32) ([]*schoolForm.UserProfile, error) {
	var faList []*schoolForm.UserProfile
	const sql = `SELECT "Profile".*, "MinimalProfile"."name", "MinimalProfile"."mobileNumber", "MinimalProfile"."phoneNumber" 
	FROM "Profile", "MinimalProfile"
	WHERE "Profile"."userId" = "MinimalProfile"."userId" AND
	"Profile"."userId" = ANY ($1) ORDER BY "userId" ASC`
	args := &pgtype.Int4Array{}
	if err := args.Set(idList); err != nil {
		log.Printf("err in setting Int4Array in Db.GetProfiles(): %v\n", err)
		return nil, errors.New("failed to setup sql for fetching profiles")
	}
	rows, err := db.DbPool.Query(ctx, sql, args)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get profiles from database: failed, %v\n", err)
		return nil, errors.New("get profiles from database: failed")
	}
	defer rows.Close()
	rs := pgxscan.NewRowScanner(rows)
	isEmpty := true
	for rows.Next() {
		var p UserProfile
		if err := rs.Scan(&p); err != nil {
			log.Printf("Scanning profile rows from db failed: %v\n", err)
			return nil, errors.New("cannot get event attendee profiles")
		}
		dto, err := p.dto()
		if err != nil {
			log.Printf("UserProfile dto error: %v\n", err)
			return nil, errors.New("cannot parse event attendee profiles")
		}
		faList = append(faList, dto)
		isEmpty = false
	}
	if isEmpty {
		log.Printf("Empty UserProfile result set\n")
		return nil, errors.New("no event attendee profiles were fetched")
	}
	return faList, nil
}

func (db *DB) GetMinProfiles(ctx context.Context, idList []int32) ([]*schoolForm.MinProfile, error) {
	var aList []*schoolForm.MinProfile
    // Allow for empty param query (empty list instead of err)
    if len(idList) == 0 {
        return aList, nil
    }
	var args = pgtype.Int4Array{}
	if err := args.Set(idList); err != nil {
		log.Printf("err in setting Int4Array in Db.GetMinProfiles(): %v\n", err)
		return nil, errors.New("failed to setup sql for fetching minimal profiles")
	}
	const sql = `SELECT * FROM "MinimalProfile" WHERE "userId" = ANY ($1) ORDER BY "userId" ASC`
	rows, err := db.DbPool.Query(ctx, sql, idList)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get minimal profiles from database: failed: %v\n", err)
		return nil, errors.New("get minimal profiles from database: failed")
	}
	defer rows.Close()
	rs := pgxscan.NewRowScanner(rows)
	isEmpty := true
	for rows.Next() {
		var m MinProfile
		if err := rs.Scan(&m); err != nil {
			log.Printf("Scanning minimal profile rows from db: failed, %v\n", err)
			return nil, errors.New("cannot get event attendee minimal profile")
		}
		aList = append(aList, m.dto())
		isEmpty = false
	}
	if isEmpty {
		log.Printf("Empty MinProfile result set\n")
		return nil, errors.New("no event attendee minimal profiles were fetched")
	}
	return aList, nil
}

func (db *DB) GetMemberByDept(ctx context.Context, des string) (*schoolForm.MinProfile, error) {
	var member MinProfile
	const sql = `SELECT * FROM "MinimalProfile"
	WHERE "userId" = ANY (SELECT "userId" FROM "Department" WHERE "description" LIKE $1)`
	rows, err := db.DbPool.Query(ctx, sql, des)
	if err == nil {
		defer rows.Close()
		if err := pgxscan.ScanOne(&member, rows); err != nil {
			log.Printf("Scanning into struct in Db.GetMemberByDept() failed: %v\n", err)
			return nil, errors.New("failed to parse member data from database")
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}
	if err != nil {
		log.Printf("Get member by department from database: failed, %v\n", err)
		return nil, errors.New("failed to get member by department from database")
	}
	return member.dto(), nil
}
