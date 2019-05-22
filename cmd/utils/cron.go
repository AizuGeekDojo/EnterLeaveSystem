package utils

import (
	"fmt"

	"gopkg.in/robfig/cron.v2"
)

func CronInit() {
	c := cron.New()
	// sec min hour date month week
	// 0 0 * * * * -> 0:0:* */* (*) Every month, day, hour and 0min, 0sec == Every hour
	c.AddFunc("0 0 * * * *", CronTestEveryHour)
	c.Start()
}

func CronTestEveryHour() {
	fmt.Println("Every hour")
}
