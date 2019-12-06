package services

import (
	"dds-backend/database"
	"dds-backend/models"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sort"
	"strconv"
	"time"
)

func launchService() {

}

func RemoveWorkerTakenItems(username string) error {
	searchItem := models.TakenItem{TakenBy: username}
	res := database.DB.Model(&models.TakenItem{}).Where(&searchItem).Delete(&searchItem)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

type ScheduleNotFoundError struct {
	error string
}

func (e *ScheduleNotFoundError) Error() string {
	return fmt.Sprintf("schedule for user not found")
}

// Extract typed schedule values from database
func GetSchedule(username string) (models.TimePoint, models.TimePoint, []models.Weekday, error) {
	searchSchedule := models.UserSchedule{Username: username}
	res := database.DB.Model(&models.UserSchedule{}).Where(&searchSchedule).First(&searchSchedule)
	if res.RecordNotFound() {
		return models.TimePoint{}, models.TimePoint{}, nil, &ScheduleNotFoundError{}
	} else if res.Error != nil {
		return models.TimePoint{}, models.TimePoint{}, nil, errors.New("server error")
	} else {
		wks, err := models.LoadWeekdays(searchSchedule.Workdays)
		if err != nil {
			return models.TimePoint{}, models.TimePoint{}, nil, err
		}
		tp1, err := models.LoadTimePoint(searchSchedule.StartTime)
		if err != nil {
			return models.TimePoint{}, models.TimePoint{}, nil, err
		}
		tp2, err := models.LoadTimePoint(searchSchedule.EndTime)
		if err != nil {
			return models.TimePoint{}, models.TimePoint{}, nil, err
		}
		return tp1, tp2, wks, nil
	}
}

func RemoveSchedule(username string) error {
	searchSchedule := models.UserSchedule{Username: username}
	res := database.DB.Model(&models.UserSchedule{}).Where(&searchSchedule).First(&searchSchedule)
	if res.RecordNotFound() {
		// do nothing // TODO is this a correct behaviour?
	} else if res.Error != nil {
		return res.Error
	} else {
		res := database.DB.Model(&models.UserSchedule{}).Delete(&searchSchedule)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

// Transform schedule into single formatted string for display in telegram bot
func PrettySchedule(startTime models.TimePoint, endTime models.TimePoint, workdays []models.Weekday) string {
	return fmt.Sprintf("Your schedule is:\n%s - %s\n%s", startTime.ToStr(),
		endTime.ToStr(), models.PrettyWeekdays(workdays))
}

// Set schedule for user, either create new or change existing
// Event is saved into database and to running cron instance
func SetSchedule(username string, gametype string, workdays []models.Weekday, startTime, endTime models.TimePoint) error {
	txn := database.DB.Begin()
	searchSchedule := models.UserSchedule{Username: username}
	res := txn.Model(&models.UserSchedule{}).Where(&searchSchedule).First(&searchSchedule)

	// if schedule exists field will be non nil
	saveSchedule := models.UserSchedule{StartCronID: searchSchedule.StartCronID, EndCronID: searchSchedule.EndCronID}

	// add cron events
	searchSchedule.StartTime = startTime.ToStr()
	searchSchedule.EndTime = endTime.ToStr()
	searchSchedule.Workdays = models.StoreWeekdays(workdays)
	searchSchedule.GameType = gametype
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
		} else {
			// on success remove old schedules from cron
			RemoveCronRange(CronInstance, saveSchedule)
		}
	}
	txn.Commit()
	return nil
}

// Format event time and date occurances into cron string
func GetCronString(time models.TimePoint, days []models.Weekday) (string, error) {
	weekformat := ""
	weekcont := []int{}
	for _, e := range models.SortSetWeekdays(days) {
		if !e.IsValid() {
			return "", errors.New("invalid weekday")
		}
		weekcont = append(weekcont, int(e)%7) // sunday is 0, saturday is 6
	}
	sort.Ints(weekcont)
	for _, e := range weekcont {
		weekformat += "," + strconv.Itoa(int(e)%7)
	}
	if len(weekformat) <= 0 {
		return "", errors.New("empty time")
	} else {
		weekformat = weekformat[1:]
	}
	res := fmt.Sprintf("%d %d * * %s", time.Minute, time.Hour, weekformat)
	return res, nil
}

// Send notification to telegram user with message
func ScheduleNotify(message string, username string) {
	err := SendNotification(username, message)
	if err != nil {
		log.Println("unable to send notification")
	}
}

// Add Start-End Time Events to Cron service
func AddCronRange(c *cron.Cron, schedule models.UserSchedule) (cron.EntryID, cron.EntryID, error) {
	wks, err := models.LoadWeekdays(schedule.Workdays)
	if err != nil {
		return 0, 0, err
	}
	timeS, err := models.LoadTimePoint(schedule.StartTime)
	if err != nil {
		return 0, 0, err
	}
	cstr, err := GetCronString(timeS, wks)
	if err != nil {
		return 0, 0, err
	}
	msgS := fmt.Sprintf("Your workday has begun (%s)", schedule.StartTime)
	id1, err := c.AddFunc(cstr, func() { ScheduleNotify(msgS, schedule.Username) })
	if err != nil {
		return 0, 0, err
	}
	timeE, err := models.LoadTimePoint(schedule.EndTime)
	if err != nil {
		c.Remove(id1)
		return 0, 0, err
	}
	cstr, err = GetCronString(timeE, wks)
	if err != nil {
		c.Remove(id1)
		return 0, 0, err
	}
	msgE := fmt.Sprintf("Your workday has finished (%s). Your taken items were removed.", schedule.EndTime)
	id2, err := c.AddFunc(cstr, func() { ScheduleNotify(msgE, schedule.Username); RemoveWorkerTakenItems(schedule.Username) })
	if err != nil {
		c.Remove(id1)
		return 0, 0, err
	}
	return id1, id2, nil
}

// Remove Start-End Time Events from Cron service
func RemoveCronRange(c *cron.Cron, schedule models.UserSchedule) {
	c.Remove(schedule.StartCronID)
	c.Remove(schedule.EndCronID)
}

// When run on start of program restores state of Cron service from database
func PerformDBCronRecovery(c *cron.Cron) {
	var schedules []models.UserSchedule
	database.DB.Model(&models.UserSchedule{}).Find(&schedules)
	for _, elem := range schedules {
		_, _, err := AddCronRange(c, elem)
		if err != nil {
			log.Printf("unable to restore cron ranges for: %s\n", elem.Username)
		}
	}
}

// Instance of Cron Scheduler
var CronInstance *cron.Cron

// Initialize Cron Scheduler service
func LaunchScheduler() {
	CronInstance = cron.New(cron.WithLocation(time.UTC))
	PerformDBCronRecovery(CronInstance)
	CronInstance.Start()
}
