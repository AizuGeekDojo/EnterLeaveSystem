package utils

import (
	"database/sql"
	"log"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
	"gopkg.in/robfig/cron.v2"
)

// CronInit setups gocron and starts cron
func CronInit(d *sql.DB) error {
	c := cron.New()
	// sec min hour date month week
	// 0 0 * * * * -> 0:0:* */* (*) Every month, day, hour and 0min, 0sec == Every hour
	_, err := c.AddFunc("0 0 0 * * *", cronForceLeave(d))
	if err != nil {
		return err
	}
	_, err = c.AddFunc("0 18 16 1 * *", cronSendMonthlyLog(d))
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func cronForceLeave(d *sql.DB) func() {
	return func() {
		err := db.ForceLeave(d)
		if err != nil {
			log.Printf("Cron: db.ForceLeave error: %v\n", err)
		}
	}
}

func cronSendMonthlyLog(d *sql.DB) func() {
	return func() {
		err := sendMonthlyLog(d)
		if err != nil {
			log.Printf("Cron: db.ForceLeave error: %v\n", err)
		} else {
			log.Printf("Monthly log exported successfully!\n")
		}
	}
}
