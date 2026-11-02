package utils

import (
	"github.com/kelseyhightower/envconfig"
)

// SlackInfo is structure for slack info in env
type SlackInfo struct {
	AppToken        string `envconfig:"APP_TOKEN" default:""`
	BotToken        string `envconfig:"BOT_TOKEN" default:""`
	UserName        string `envconfig:"USER_NAME" default:"Logging Bot"`
	IconEmoji       string `envconfig:"ICON_EMOJI" default:":robot_face:"`
	ChannelID       string `envconfig:"CHANNEL_ID" default:"#enter_leave_log"`
	ErrorChannelID  string `envconfig:"ERROR_CHANNEL_ID" default:"#_sa_system"`
	CSVLogToken     string `envconfig:"CSV_LOG_TOKEN" default:""`
	CSVLogChannelID string `envconfig:"CSV_LOG_CHANNEL_ID" default:""`
}

// ReaderInfo is structure for card reader info in env
type ReaderInfo struct {
	// Mock reads cards from stdin instead of the RC-S380 (for development)
	Mock bool `envconfig:"NFC_MOCK" default:"false"`
}

var (
	slackinfo  SlackInfo
	readerinfo ReaderInfo
)

func init() {
	if err := envconfig.Process("AGD", &slackinfo); err != nil {
		panic(err)
	}
	if err := envconfig.Process("AGD", &readerinfo); err != nil {
		panic(err)
	}
}

// GetSlackInfo returns slack config info
func GetSlackInfo() SlackInfo {
	return slackinfo
}

// GetReaderInfo returns card reader config info
func GetReaderInfo() ReaderInfo {
	return readerinfo
}
