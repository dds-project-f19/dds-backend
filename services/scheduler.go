package services

import (
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func launchService() {

}

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

func SetSchedule(username string, workdays []Weekday, startTime, endTime time.Time) {
	// TODO search for existing schedule
	// TODO delete existing schedule if exists by id
	// TODO create new schedule and save id
}

func LaunchScheduler() {
	c := cron.New() // TODO schare cron link
	// TODO recover existing schedules from database
	_, err := c.AddFunc("20 21 * * 1-5", func() { fmt.Println("At 21:18 monday to friday") })
	if err != nil {
		fmt.Println("CRON ERROR")
	}
	c.Start()

}
