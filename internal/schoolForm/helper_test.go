package schoolForm

import (
	"fmt"
	"strings"
	"time"
)

func getFullEventInfo() *EventInfo {
	return &EventInfo{
		Id:             1,
		InviteToken:    "zVztX7II4eJi9b0OrV5Zj",
		Title:          "Event #1",
		BeginDate:      time.Date(2022, 8, 23, 8, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2022, 8, 28, 8, 0, 0, 0, time.UTC),
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
		EquipList: []Equip{
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
		TechEquipList: []Equip{
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
		Attendants: []FullAttendance{
			{
				UserId: 1,
				Role:   "Host",
				Jobs:   "領隊、證保",
				UserProfile: UserProfile{
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
			},
			{
				UserId: 2,
				Role:   "Mentor",
				Jobs:   "輔隊",
				UserProfile: UserProfile{
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
			},
			{
				UserId: 3,
				Role:   "Member",
				Jobs:   "大廚、裝備、學員",
				UserProfile: UserProfile{
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
			},
			{
				UserId: 4,
				Role:   "Member",
				Jobs:   "",
				UserProfile: UserProfile{
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
		Rescues: []Attendance{
			{
				UserId: 5,
				Role:   "Rescue",
				MinProfile: MinProfile{
					UserId:       5,
					Name:         "半號一君",
					MobileNumber: "0910-000-004",
					PhoneNumber:  "01-0000004",
				},
			},
		},
		Watchers: []Attendance{
			{
				UserId: 6,
				Role:   "Watcher",
				MinProfile: MinProfile{
					UserId:       6,
					Name:         "半號二君",
					MobileNumber: "0910-000-005",
					PhoneNumber:  "01-0000005",
				},
			},
			{
				UserId: 7,
				Role:   "Watcher",
				MinProfile: MinProfile{
					UserId:       7,
					Name:         "半號三君",
					MobileNumber: "0910-000-006",
					PhoneNumber:  "01-0000006",
				},
			},
		},
	}
}

func getEInfo_longFields() *EventInfo {
	baseEInfo := getFullEventInfo()
	longText := []string{
		"Nunc dapibus ut tellus nec viverra. In dapibus ex sit amet mauris aliquam, vel luctus tellus commodo. Morbi in leo id erat pretium tempor. Sed bibendum dui lacus, at vulputate dui cursus ut.",
		"Maecenas nec commodo diam. Interdum et malesuada fames ac ante ipsum primis in faucibus. Quisque dictum sapien gravida arcu accumsan, ac hendrerit turpis aliquet. Mauris eleifend sem sodales ipsum rutrum pharetra. Nam ullamcorper condimentum mi, id laoreet nisl aliquam sit amet.",
		"Nunc pulvinar ante quis justo suscipit, in luctus quam sodales. Morbi commodo erat et libero mattis porta. Duis tellus eros, fermentum ut commodo non, blandit non neque. Suspendisse vitae est id ex aliquam varius. Aliquam facilisis urna ut lobortis tincidunt. Nullam a nisi vitae sem posuere placerat quis ac odio. Morbi sed eros ante. Sed euismod faucibus neque, quis ullamcorper est pretium id.",
		"Quisque at est mollis, interdum leo non, ultricies ligula. Phasellus ut augue semper, mollis lacus eu, blandit est. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Pellentesque commodo lacus nec elementum efficitur. Mauris varius faucibus diam, sit amet consequat leo tempus ut. Maecenas et ligula congue, cursus ipsum non, ultricies augue.",
	}

	/* Writing shorter texts */
	// #of sentences in the shorter text
	shortSentenceLen := 2
	var sB strings.Builder
	for i := 0; i < shortSentenceLen; i++ {
		if _, err := sB.WriteString(longText[i] + "\n"); err != nil {
			// todo: i'm just outputting err string here, i should do something about uncaught err?
			fmt.Println(err.Error())
		}
	}
	baseEInfo.RetreatPlan = sB.String()

	/* Writing longer texts */
	for i := shortSentenceLen; i < len(longText); i++ {
		if _, err := sB.WriteString(longText[i] + "\n"); err != nil {
			// todo: i'm just outputting err string here, i should do something about uncaught err?
			fmt.Println(err.Error())
		}
	}
	baseEInfo.TripOverview = sB.String()
	baseEInfo.Records = sB.String()

	/* Writing to equip section */
	baseEInfo.EquipList = []Equip{
		{Name: "帳棚", Des: "1x"},
		{Name: "鍋組（含湯瓢、鍋夾）", Des: "一中、兩大、二小（加兩個勺子）(I need some more words here XDDDDD)"},
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
		{Name: "Fibriophobia(Having fear of fever)", Des: "1x"},
		{Name: "Utilitarianism", Des: "(Adopting a code of conduct that determines ethical values)"},
		{Name: "ooxx", Des: "1x"},
		{Name: "xxoo", Des: "1x"},
		{Name: "ooxx", Des: "1x"},
		{Name: "xxoo", Des: "1x"},
	}
	baseEInfo.TechEquipList = []Equip{
		{Name: "主繩", Des: "1x"},
		{Name: "吊帶", Des: "2x"},
		{Name: "上升器", Des: "2x"},
		{Name: "下降器", Des: "2x"},
		{Name: "岩盔", Des: "2x"},
		{Name: "有鎖鉤環", Des: "4x"},
		{Name: "救生衣", Des: "4x"},
		{Name: "岩釘", Des: "Angle*1, Knifeblade*5, Arrow*3, Beak*10, Skyhook*10"},
		{Name: "ooxx", Des: "1x"},
		{Name: "oxxo", Des: "1x"},
		{Name: "oxox", Des: "1x"},
	}

	return baseEInfo
}
