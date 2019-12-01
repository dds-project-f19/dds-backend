package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"sort"
	"strconv"
	"strings"
)

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (w *Weekday) IsValid() bool {
	return int(*w) >= 1 && int(*w) <= 7
}

func (w *Weekday) ToStr() string {
	if !w.IsValid() {
		return "ErrorWeekday"
	}
	return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}[int(*w)-1]
}

func SortSetWeekdays(days []Weekday) []Weekday {
	daysMap := map[int]bool{}
	var outDays []Weekday
	var intDays []int
	for _, e := range days {
		intDays = append(intDays, int(e))
	}
	sort.Ints(intDays)
	for _, e := range intDays {
		_, ok := daysMap[e]
		if !ok { // weekday not already added
			daysMap[e] = true
			outDays = append(outDays, Weekday(e))
		}
	}
	return outDays
}

func PrettyWeekdays(days []Weekday) string {
	var weekformat string
	for _, w := range SortSetWeekdays(days) {
		weekformat += w.ToStr() + " "
	}
	return strings.TrimSpace(weekformat)
}

func StoreWeekdays(days []Weekday) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(SortSetWeekdays(days))), ","), "[]")
}

func LoadWeekdays(in string) ([]Weekday, error) {
	in = strings.TrimSpace(in)
	if in == "" {
		return []Weekday{}, nil
	}
	vals := strings.Split(in, ",")
	var store []Weekday
	for _, elem := range vals {
		i, err := strconv.Atoi(elem)
		if err != nil {
			return nil, err
		}
		wk := Weekday(i)
		if !wk.IsValid() {
			return nil, errors.New("unexpected weekday index")
		}
		store = append(store, Weekday(i))
	}
	return SortSetWeekdays(store), nil
}

type TimePoint struct {
	Hour   int
	Minute int
}

func LoadTimePoint(in string) (TimePoint, error) {
	tp := TimePoint{}
	read, err := fmt.Sscanf(in, "%d:%d", &tp.Hour, &tp.Minute)
	if err != nil {
		return TimePoint{}, err
	}
	if read != 2 {
		return TimePoint{}, errors.New("bad number format")
	}
	return tp, nil
}

func (t *TimePoint) ToStr() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

func (t1 *TimePoint) Before(t2 TimePoint) bool {
	return t1.Hour < t2.Hour || (t1.Hour == t2.Hour) && (t1.Minute < t2.Minute)
}

func (t *TimePoint) IsValid() bool {
	return 23 >= t.Hour && t.Hour >= 0 && 59 >= t.Minute && t.Minute >= 0
}

type UserSchedule struct {
	gorm.Model
	Username    string       `gorm:"unique_index;not null"`
	StartTime   string       `gorm:"not null"`
	EndTime     string       `gorm:"not null"`
	Workdays    string       `gorm:"not null"`
	StartCronID cron.EntryID `gorm:"not null"`
	EndCronID   cron.EntryID `gorm:"not null"`
	GameType    string       `gorm:"not null"`
}
