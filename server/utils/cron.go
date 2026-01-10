package utils

import (
	"database/sql"
	"log"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"github.com/robfig/cron/v3"
)

// CronInit setups cron scheduler and starts cron jobs
func CronInit(d *sql.DB) error {
	// Create cron with second precision
	c := cron.New(cron.WithSeconds())

	// Force leave all users at midnight every day
	// Format: sec min hour day month day-of-week
	// "0 0 0 * * *" = Every day at 00:00:00
	_, err := c.AddFunc("0 0 0 * * *", cronForceLeave(d))
	if err != nil {
		return err
	}

	// Send monthly log on the 1st of every month at 01:00:00
	// "0 0 1 1 * *" = 1st day of month at 01:00:00
	_, err = c.AddFunc("0 0 1 1 * *", cronSendMonthlyLog(d))
	if err != nil {
		return err
	}

	c.Start()
	log.Println("Cron scheduler started successfully")
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
