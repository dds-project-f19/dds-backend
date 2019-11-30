package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
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

func (w *Weekday) ToStr() string {
	return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}[int(*w)]
}

func PrettyWeekdays(days []Weekday) string {
	weekdays := map[int]bool{}
	var weekformat string
	var intDays []int
	for _, e := range days {
		intDays = append(intDays, int(e))
	}
	sort.Ints(intDays)
	for _, e := range intDays {
		_, ok := weekdays[e]
		if !ok { // weekday not already added
			weekdays[e] = true
			wk := Weekday(e)
			weekformat += wk.ToStr() + " "
		}
	}
	return weekformat
}

func StoreWeekdays(days []Weekday) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(days)), ","), "[]")
}

func LoadWeekdays(in string) ([]Weekday, error) {
	vals := strings.Split(in, ",")
	var store []Weekday
	for _, elem := range vals {
		i, err := strconv.Atoi(elem)
		if err != nil {
			return nil, err
		}
		store = append(store, Weekday(i))
	}
	return store, nil
}

type TimePoint struct {
	Hour   int
	Minute int
}

func (t1 *TimePoint) Before(t2 TimePoint) bool {
	return t1.Hour < t2.Hour || (t1.Hour == t2.Hour) && (t1.Minute < t2.Minute)
}

func (t *TimePoint) IsValid() bool {
	return 23 >= t.Hour && t.Hour >= 0 && 59 >= t.Minute && t.Minute >= 0
}

func (t *TimePoint) ToStr() string {
	return fmt.Sprintf("%d:%d", t.Hour, t.Minute)
}

type UserSchedule struct {
	gorm.Model
	Username    string       `gorm:"unique_index;not null"`
	StartTime   TimePoint    `gorm:"not null"`
	EndTime     TimePoint    `gorm:"not null"`
	Workdays    string       `gorm:"not null"`
	StartCronID cron.EntryID `gorm:"not null"`
	EndCronID   cron.EntryID `gorm:"not null"`
}
