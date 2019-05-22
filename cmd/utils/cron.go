package utils

import (
	"database/sql"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
	"gopkg.in/robfig/cron.v2"
)

func CronInit(d *sql.DB) {
	c := cron.New()
	// sec min hour date month week
	// 0 0 * * * * -> 0:0:* */* (*) Every month, day, hour and 0min, 0sec == Every hour
	c.AddFunc("0 0 0 * * *", func() {
		db.ForceLeave(d)
	})
	c.Start()
}
