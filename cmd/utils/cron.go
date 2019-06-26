package utils

import (
	"database/sql"
	"log"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
	"gopkg.in/robfig/cron.v2"
)

func CronInit(d *sql.DB) error {
	c := cron.New()
	// sec min hour date month week
	// 0 0 * * * * -> 0:0:* */* (*) Every month, day, hour and 0min, 0sec == Every hour
	_, err := c.AddFunc("0 0 0 * * *", func() {
		err := db.ForceLeave(d)
		if err != nil {
			log.Printf("Cron: db.ForceLeave error: %v", err)
		}
	})
	if err != nil {
		return err
	}
	_, err = c.AddFunc("0 0 0 1 * *", func() {
		err := SendMonthlyLog(d)
		if err != nil {
			log.Printf("Cron: db.ForceLeave error: %v", err)
		}
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
