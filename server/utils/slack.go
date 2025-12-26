package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var (
	slackClient     *slack.Client
	socketClient    *socketmode.Client
	slackClientOnce sync.Once
)

// InitSlackSocketMode initializes the Slack Socket Mode client
func InitSlackSocketMode() error {
	var err error
	slackClientOnce.Do(func() {
		cfg := GetSlackInfo()
		if cfg.BotToken == "" {
			err = errors.New("Slack Bot Token is not defined")
			return
		}

		slackClient = slack.New(
			cfg.BotToken,
			slack.OptionAppLevelToken(cfg.AppToken),
		)

		if cfg.AppToken != "" {
			socketClient = socketmode.New(
				slackClient,
				socketmode.OptionDebug(false),
			)
			log.Println("Slack Socket Mode initialized")
		} else {
			log.Println("Slack client initialized (no App Token for Socket Mode)")
		}
	})
	return err
}

// GetSlackClient returns the initialized Slack client
func GetSlackClient() *slack.Client {
	if slackClient == nil {
		if err := InitSlackSocketMode(); err != nil {
			log.Printf("Failed to initialize Slack client: %v", err)
			return nil
		}
	}
	return slackClient
}

// SlackNotify sends Slack notification using slack-go socket mode client.
func SlackNotify(name string, uid string, isEnter bool, timestamp time.Time, ext string) error {
	client := GetSlackClient()
	if client == nil {
		return errors.New("Slack client is not initialized")
	}

	cfg := GetSlackInfo()
	if cfg.ChannelID == "" {
		return errors.New("Slack ChannelID is not defined")
	}

	status := "退室"
	if isEnter {
		status = "入室"
	}

	messageText := fmt.Sprintf("%v : %v さんが %v に%vしました。", uid, name, timestamp.Format("2006-01-02 15:04:05"), status)

	var attachments []slack.Attachment

	if ext != "" {
		var extMap map[string]interface{}
		if err := json.Unmarshal([]byte(ext), &extMap); err != nil {
			log.Println(err)
			log.Println("Error has occured in SlackNotify. User Ext is : ", ext)
			return err
		}

		usage := extMap["Use"]
		mess := extMap["message"]
		attachments = append(attachments, slack.Attachment{
			Title: "アンケート結果",
			Text:  fmt.Sprintf("目的 : %v \n 感想 : %v", usage, mess),
		})
	}
	log.Println(cfg.ChannelID)
	_, _, err := client.PostMessage(
		cfg.ChannelID,
		slack.MsgOptionText(messageText, false),
		slack.MsgOptionUsername(cfg.UserName),
		slack.MsgOptionIconEmoji(cfg.IconEmoji),
		slack.MsgOptionAttachments(attachments...),
	)

	if err != nil {
		log.Printf("Failed to send Slack message: %v", err)
		return err
	}

	return nil
}
