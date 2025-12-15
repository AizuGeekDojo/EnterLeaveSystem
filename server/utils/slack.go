package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/slack-go/slack"
)

// SlackNotify sends Slack notification using slack-go webhook client.
func SlackNotify(name string, uid string, isEnter bool, timestamp time.Time, ext string) error {
	cfg := GetSlackInfo()
	webhookURL := cfg.GetWebHookURL()
	if webhookURL == "" {
		return errors.New("Slack URL is not defined")
	}

	status := "退室"
	if isEnter {
		status = "入室"
	}

	message := &slack.WebhookMessage{
		Channel:   cfg.Channel,
		Username:  cfg.UserName,
		IconEmoji: cfg.IconEmoji,
		Text:      fmt.Sprintf("%v : %v さんが %v に%vしました。", uid, name, timestamp.Format("2006-01-02 15:04:05"), status),
	}

	if ext != "" {
		var extMap map[string]interface{}
		if err := json.Unmarshal([]byte(ext), &extMap); err != nil {
			log.Println(err)
			log.Println("Error has occured in SlackNotify. User Ext is : ", ext)
			return err
		}

		usage := extMap["Use"]
		mess := extMap["message"]
		message.Attachments = append(message.Attachments, slack.Attachment{
			Title: "アンケート結果",
			Text:  fmt.Sprintf("目的 : %v \n 感想 : %v", usage, mess),
		})
	}

	if err := slack.PostWebhook(webhookURL, message); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
