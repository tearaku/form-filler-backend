package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"teacup1592/form-filler/internal/dataSrc"
	"teacup1592/form-filler/internal/schoolForm"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	Db  *DB
	ctx context.Context
}

func (s *DBTestSuite) SetupSuite() {
	if err := dataSrc.LocalEnvSetupInTest(); err != nil {
		s.T().Fatal(err.Error())
	}

	testDbURL := "postgres://postgres:admin@localhost:5433/postgres"
	connPool, err := NewDBPool(context.Background(), testDbURL)
	if err != nil {
		s.T().Fatal(err.Error())
	}
	s.Db = &DB{DbPool: connPool}
	s.ctx = context.Background()

	// TODO: can optimize it using Batch instead of individual queries
	// Optionally generate data in DB via env var
	if genData, err := strconv.Atoi(os.Getenv("GEN_TEST_DATA")); err != nil || genData == 0 {
		return
	}

	if err := s.Db.DbPool.BeginFunc(s.ctx, func(tx pgx.Tx) error {
		fmt.Println("Generating mock data to dockerized Postgres db...")
		if err := createUsers(s.ctx, tx); err != nil {
			return err
		}
		if err := createEvent(s.ctx, tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		s.T().Fatal(err.Error())
	}
}

func (s *DBTestSuite) TearDownSuite() {
	s.Db.DbPool.Close()
}

func (s *DBTestSuite) TestGetEventInfo() {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		want    *schoolForm.EventInfo
		wantErr string
	}{
		{
			name: "get valid event",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &schoolForm.StubEventInfoData,
		},
		{
			name: "get event with non-existent id",
			args: args{
				ctx: context.Background(),
				id:  9999,
			},
			wantErr: "failed to get event from database",
		},
	}
	// These two fields are only partially filled with this db service call
	cmpOpt := []cmp.Option{
		cmpopts.IgnoreFields(schoolForm.FullAttendance{}, "UserProfile"),
		cmpopts.IgnoreFields(schoolForm.Attendance{}, "MinProfile"),
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.Db.GetEventInfo(tt.args.ctx, tt.args.id)
			if err == nil && tt.wantErr != "" || err != nil && err.Error() != tt.wantErr {
				t.Errorf("Db.GetEventInfo() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want, cmpOpt...) {
				t.Errorf("mismatch with value from Db.GetEventInfo() [-got, +want]: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func (s *DBTestSuite) TestGetProfiles() {
	type args struct {
		ctx    context.Context
		idList []int32
	}
	tests := []struct {
		name    string
		args    args
		want    []*schoolForm.UserProfile
		wantErr string
	}{
		{
			name: "get valid profiles",
			args: args{
				ctx:    context.Background(),
				idList: []int32{1, 2, 3, 4},
			},
			want: []*schoolForm.UserProfile{
				&schoolForm.StubProfileList[0],
				&schoolForm.StubProfileList[1],
				&schoolForm.StubProfileList[2],
				&schoolForm.StubProfileList[3],
			},
		},
		{
			name: "get users without profiles (Profile)",
			args: args{
				ctx:    context.Background(),
				idList: []int32{5, 6, 7},
			},
			wantErr: "no event attendee profiles were fetched",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.Db.GetProfiles(tt.args.ctx, tt.args.idList)
			if err == nil && tt.wantErr != "" || err != nil && err.Error() != tt.wantErr {
				t.Errorf("Db.GetProfiles() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mismatch with value from Db.GetProfiles(): %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func (s *DBTestSuite) TestGetMinProfiles() {
	type args struct {
		ctx    context.Context
		idList []int32
	}
	tests := []struct {
		name    string
		args    args
		want    []*schoolForm.MinProfile
		wantErr string
	}{
		{
			name: "get valid min profiles",
			args: args{
				ctx:    context.Background(),
				idList: []int32{1, 2, 3, 4, 5, 6, 7},
			},
			want: []*schoolForm.MinProfile{
				&schoolForm.StubMinProfileList[0],
				&schoolForm.StubMinProfileList[1],
				&schoolForm.StubMinProfileList[2],
				&schoolForm.StubMinProfileList[3],
				&schoolForm.StubMinProfileList[4],
				&schoolForm.StubMinProfileList[5],
				&schoolForm.StubMinProfileList[6],
			},
		},
		{
			name: "get with empty id list",
			args: args{
				ctx:    context.Background(),
				idList: []int32{},
			},
			want: nil,
		},
		{
			name: "get non-existent min profiles",
			args: args{
				ctx:    context.Background(),
				idList: []int32{1000, 1001, 1002, 1003},
			},
			wantErr: "no event attendee minimal profiles were fetched",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.Db.GetMinProfiles(tt.args.ctx, tt.args.idList)
			if err == nil && tt.wantErr != "" || err != nil && err.Error() != tt.wantErr {
				t.Errorf("Db.GetMinProfiles() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mismatch with value from Db.GetMinProfiles(): %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func (s *DBTestSuite) TestGetMemberByDept() {
	type args struct {
		ctx context.Context
		des string
	}
	tests := []struct {
		name    string
		args    args
		want    *schoolForm.MinProfile
		wantErr string
	}{
		{
			name: "get valid departmenet head: 社長",
			args: args{
				ctx: context.Background(),
				des: "社長",
			},
			want: &schoolForm.StubMinProfileList[0],
		},
		{
			name: "get valid departmenet head: 嚮導部長1",
			args: args{
				ctx: context.Background(),
				des: "嚮導部長1",
			},
			want: &schoolForm.StubMinProfileList[1],
		},
		{
			name: "get valid departmenet head: 嚮導部長2",
			args: args{
				ctx: context.Background(),
				des: "嚮導部長2",
			},
			want: &schoolForm.StubMinProfileList[2],
		},
		{
			name: "get valid multi-hit on department head: 嚮導部長",
			args: args{
				ctx: context.Background(),
				des: "嚮導部長",
			},
			want: &schoolForm.StubMinProfileList[1],
		},
		{
			name: "get valid departmenet head: 社產組長",
			args: args{
				ctx: context.Background(),
				des: "社產組長",
			},
			want: &schoolForm.StubMinProfileList[3],
		},
		{
			name: "get valid departmenet head: 山難部長",
			args: args{
				ctx: context.Background(),
				des: "山難部長",
			},
			want: &schoolForm.StubMinProfileList[4],
		},
		{
			name: "get non-existent department head",
			args: args{
				ctx: context.Background(),
				des: "香菇部長",
			},
			wantErr: "no members with specified department is fetched",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.Db.GetMemberByDept(tt.args.ctx, tt.args.des)
			if err == nil && tt.wantErr != "" || err != nil && err.Error() != tt.wantErr {
				t.Errorf("Db.GetMemberByDept() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mismatch with value from Db.GetMemberByDept(): %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
