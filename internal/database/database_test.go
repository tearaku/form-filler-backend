package database

import (
	"context"
	"os"
	"teacup1592/form-filler/internal/dataSrc"
	"teacup1592/form-filler/internal/schoolForm"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	Db *DB
}

func (s *DBTestSuite) SetupSuite() {
	if err := dataSrc.LocalEnvSetupInTest(); err != nil {
		s.T().Fatal(err.Error())
	}
	connPool, err := NewDBPool(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		s.T().Fatal(err.Error())
	}
	s.Db = &DB{DbPool: connPool}
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
			want: &schoolForm.EventInfo{
				Id:             1,
				InviteToken:    "zVztX7II4eJi9b0OrV5Zj",
				Title:          "Event #1",
				BeginDate:      time.Date(2022, 7, 23, 16, 0, 0, 0, time.UTC),
				EndDate:        time.Date(2022, 7, 28, 16, 0, 0, 0, time.UTC),
				Location:       "Taipei",
				Category:       "B勘",
				GroupCategory:  "天狼",
				Drivers:        "司機一號、司機二號",
				DriversNumber:  "0900-111-111, 0900-111-112",
				RadioFreq:      "145.20 Mhz",
				RadioCodename:  "浩浩",
				TripOverview:   "D0 wwwwwww\nD1 oooooooo\nD2 zzzzzzzzz\nD3 qqqqqq",
				RescueTime:     "D5 1800",
				RetreatPlan:    "C3 沒過ＯＯＸＸ，原路哈哈哈",
				MapCoordSystem: "TWD97 上河",
				Records:        "[0] ooxx/oo/xx wwoowwoo\n[1] ooxx/xx/oo oxoxoxoxox\n",
				EquipList: []schoolForm.Equip{
					{Name: "帳棚", Des: "1x"},
					{Name: "鍋組（含湯瓢、鍋夾）", Des: "1x"},
					{Name: "爐頭", Des: "1x"},
					{Name: "Gas", Des: "1x"},
					{Name: "糧食", Des: "1x"},
					{Name: "預備糧", Des: "1x"},
					{Name: "山刀", Des: "1x"},
					{Name: "鋸子", Des: "1x"},
					{Name: "路標", Des: "1x"},
					{Name: "衛星電話", Des: "1x"},
					{Name: "收音機", Des: "1x"},
					{Name: "無線電", Des: "1x"},
					{Name: "傘帶", Des: "1x"},
					{Name: "Sling", Des: "1x"},
					{Name: "無鎖鉤環", Des: "1x"},
					{Name: "急救包", Des: "1x"},
					{Name: "GPS", Des: "1x"},
					{Name: "包溫瓶", Des: "1x"},
					{Name: "ooxx", Des: "1x"},
					{Name: "xxoo", Des: "1x"},
					{Name: "ooxx", Des: "1x"},
					{Name: "xxoo", Des: "1x"},
					{Name: "ooxx", Des: "1x"},
					{Name: "xxoo", Des: "1x"},
				},
				TechEquipList: []schoolForm.Equip{
					{Name: "主繩", Des: "1x"},
					{Name: "吊帶", Des: "2x"},
					{Name: "上升器", Des: "2x"},
					{Name: "下降器", Des: "2x"},
					{Name: "岩盔", Des: "2x"},
					{Name: "有鎖鉤環", Des: "4x"},
					{Name: "救生衣", Des: "4x"},
					{Name: "ooxx", Des: "1x"},
					{Name: "ooxx", Des: "1x"},
					{Name: "oxxo", Des: "1x"},
					{Name: "oxox", Des: "1x"},
				},
				Attendants: []schoolForm.FullAttendance{
					{
						UserId: 1,
						Role:   "Host",
						Jobs:   "領隊、證保",
					},
					{
						UserId: 2,
						Role:   "Mentor",
						Jobs:   "輔隊",
					},
					{
						UserId: 3,
						Role:   "Member",
						Jobs:   "大廚、裝備、學員",
					},
					{
						UserId: 4,
						Role:   "Member",
						Jobs:   "",
					},
				},
				Rescues: []schoolForm.Attendance{
					{
						UserId: 5,
						Role:   "Rescue",
					},
				},
				Watchers: []schoolForm.Attendance{
					{
						UserId: 6,
						Role:   "Watcher",
					},
					{
						UserId: 7,
						Role:   "Watcher",
					},
				},
			},
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
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := s.Db.GetEventInfo(tt.args.ctx, tt.args.id)
			if err == nil && tt.wantErr != "" || err != nil && err.Error() != tt.wantErr {
				t.Errorf("Db.GetEventInfo() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("mismatch with value from Db.GetEventInfo(): %v", cmp.Diff(got, tt.want))
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
				{
					UserId:                 1,
					EngName:                "",
					IsMale:                 true,
					IsStudent:              true,
					MajorYear:              "昆蟲四",
					DateOfBirth:            time.Date(2000, 2, 4, 16, 0, 0, 0, time.UTC),
					PlaceOfBirth:           "呵呵市",
					IsTaiwanese:            true,
					NationalId:             "A12345678",
					PassportNumber:         "",
					Nationality:            "",
					Address:                "呵呵地址",
					EmergencyContactName:   "緊急一",
					EmergencyContactMobile: "0900-000-000",
					EmergencyContactPhone:  "04-0000000",
					BeneficiaryName:        "受益一",
					BeneficiaryRelation:    "母子",
					RiceAmount:             4,
					FoodPreference:         "喜歡辣",
					Name:                   "一號君",
					MobileNumber:           "0910-000-000",
					PhoneNumber:            "01-0000000",
				},
				{
					UserId:                 2,
					EngName:                "",
					IsMale:                 false,
					IsStudent:              true,
					MajorYear:              "中文一",
					DateOfBirth:            time.Date(2004, 2, 4, 16, 0, 0, 0, time.UTC),
					PlaceOfBirth:           "呵呵市",
					IsTaiwanese:            true,
					NationalId:             "A12345678",
					PassportNumber:         "",
					Nationality:            "",
					Address:                "呵呵地址",
					EmergencyContactName:   "緊急二",
					EmergencyContactMobile: "0900-000-001",
					EmergencyContactPhone:  "無",
					BeneficiaryName:        "受益二",
					BeneficiaryRelation:    "父女",
					RiceAmount:             2,
					FoodPreference:         "",
					Name:                   "二號君",
					MobileNumber:           "0910-000-001",
					PhoneNumber:            "01-0000001",
				},
				{
					UserId:                 3,
					EngName:                "Matthews Brittney",
					IsMale:                 true,
					IsStudent:              false,
					MajorYear:              "",
					DateOfBirth:            time.Date(1995, 2, 4, 16, 0, 0, 0, time.UTC),
					PlaceOfBirth:           "呵呵國",
					IsTaiwanese:            false,
					NationalId:             "",
					PassportNumber:         "P12345678",
					Nationality:            "香港",
					Address:                "呵呵地址",
					EmergencyContactName:   "緊急三",
					EmergencyContactMobile: "0900-000-002",
					EmergencyContactPhone:  "04-0000002",
					BeneficiaryName:        "受益三",
					BeneficiaryRelation:    "父子",
					RiceAmount:             6,
					FoodPreference:         "飯多一點",
					Name:                   "三號君",
					MobileNumber:           "0910-000-002",
					PhoneNumber:            "01-0000002",
				},
				{
					UserId:                 4,
					EngName:                "Cole Schriber",
					IsMale:                 false,
					IsStudent:              false,
					MajorYear:              "",
					DateOfBirth:            time.Date(1990, 2, 4, 16, 0, 0, 0, time.UTC),
					PlaceOfBirth:           "呵呵國",
					IsTaiwanese:            false,
					NationalId:             "",
					PassportNumber:         "P87654321",
					Nationality:            "奧門",
					Address:                "呵呵地址",
					EmergencyContactName:   "緊急四",
					EmergencyContactMobile: "0900-000-003",
					EmergencyContactPhone:  "04-0000003",
					BeneficiaryName:        "受益四",
					BeneficiaryRelation:    "母女",
					RiceAmount:             3,
					FoodPreference:         "",
					Name:                   "四號君",
					MobileNumber:           "0910-000-003",
					PhoneNumber:            "01-0000003",
				},
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

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
