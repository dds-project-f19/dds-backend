package services

import (
	"dds-backend/database"
	"dds-backend/models"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"sort"
	"strconv"
)

func launchService() {

}

// TODO Get schedule as raw values

func GetSchedulePretty(username string) (string, error) {
	searchSchedule := models.UserSchedule{Username: username}
	res := database.DB.Model(&models.UserSchedule{}).Where(&searchSchedule).First(&searchSchedule)
	if res.RecordNotFound() {
		return "", errors.New("schedule for this user not found")
	} else if res.Error != nil {
		return "", errors.New("server error")
	} else {
		wks, err := models.LoadWeekdays(searchSchedule.Workdays)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Your schedule is:\n%s - %s\n%s", searchSchedule.StartTime.ToStr(),
			searchSchedule.EndTime.ToStr(), models.PrettyWeekdays(wks)), nil
	}
}

func SetSchedule(username string, workdays []models.Weekday, startTime, endTime models.TimePoint) error {
	txn := database.DB.Begin()
	searchSchedule := models.UserSchedule{Username: username}
	res := txn.Model(&models.UserSchedule{}).Where(&searchSchedule).First(&searchSchedule)

	// add cron events
	searchSchedule.StartTime = startTime
	searchSchedule.EndTime = endTime
	searchSchedule.Workdays = models.StoreWeekdays(workdays)
	id1, id2, err := AddCronRange(CronInstance, searchSchedule)
	if err != nil {
		return err
	}
	// create new record with ids
	searchSchedule.StartCronID = id1
	searchSchedule.EndCronID = id2

	if res.RecordNotFound() {
		res := txn.Model(&models.UserSchedule{}).Create(&searchSchedule)
		if res.Error != nil {
			RemoveCronRange(CronInstance, searchSchedule)
			return res.Error
		}
	} else if res.Error != nil {
		RemoveCronRange(CronInstance, searchSchedule)
		return res.Error
	} else {
		res := txn.Model(&models.UserSchedule{}).Save(&searchSchedule)
		if res.Error != nil {
			RemoveCronRange(CronInstance, searchSchedule)
			return res.Error
		}
	}
	txn.Commit()
	return nil
}

func GetCronString(time models.TimePoint, days []models.Weekday) (string, error) {
	weekdays := map[int]bool{}
	weekformat := ""
	var intDays []int
	for _, e := range days { // TODO absolutely terrible code, refactor
		intDays = append(intDays, int(e))
	}
	sort.Ints(intDays)
	for _, e := range intDays {
		_, ok := weekdays[int(e)]
		if !ok { // weekday not already added
			weekdays[int(e)] = true
			weekformat += "," + strconv.Itoa(int(e))
		}
	}
	if len(weekformat) <= 0 {
		return "", errors.New("Empty time")
	} else {
		weekformat = weekformat[1:]
	}
	res := fmt.Sprintf("%d %d * * %s", time.Minute, time.Hour, weekformat)
	return res, nil
}

func AddCronRange(c *cron.Cron, schedule models.UserSchedule) (cron.EntryID, cron.EntryID, error) {
	wks, err := models.LoadWeekdays(schedule.Workdays)
	if err != nil {
		return 0, 0, err
	}
	cstr, err := GetCronString(schedule.StartTime, wks)
	if err != nil {
		return 0, 0, err
	}
	id1, err := c.AddFunc(cstr, func() { ScheduleNotify("Your workday has begun!", schedule.Username) })
	if err != nil {
		return 0, 0, err
	}
	cstr, err = GetCronString(schedule.EndTime, wks)
	if err != nil {
		c.Remove(id1)
		return 0, 0, err
	}
	id2, err := c.AddFunc(cstr, func() { ScheduleNotify("Your workday has finished!", schedule.Username) })
	if err != nil {
		c.Remove(id1)
		return 0, 0, err
	}
	return id1, id2, nil
}

func ScheduleNotify(message string, username string) {
	err := SendNotification(username, message)
	if err != nil {
		log.Println("unable to send notification")
	}
}

func RemoveCronRange(c *cron.Cron, schedule models.UserSchedule) {
	c.Remove(schedule.StartCronID)
	c.Remove(schedule.EndCronID)
}

func PerformDBCronRecovery(c *cron.Cron) {
	var schedules []models.UserSchedule
	database.DB.Model(&models.UserSchedule{}).Find(&schedules)
	for _, elem := range schedules {
		AddCronRange(c, elem)
	}
}

var CronInstance *cron.Cron

func LaunchScheduler() {
	CronInstance = cron.New()
	PerformDBCronRecovery(CronInstance)
	CronInstance.Start()

}
