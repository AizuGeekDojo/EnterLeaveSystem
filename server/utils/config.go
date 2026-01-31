package utils

import (
	"github.com/kelseyhightower/envconfig"
)

// SlackInfo is structure for slack info in env
type SlackInfo struct {
	AppToken    string `envconfig:"APP_TOKEN" default:""`
	BotToken    string `envconfig:"BOT_TOKEN" default:""`
	UserName    string `envconfig:"USER_NAME" default:"Logging Bot"`
	IconEmoji   string `envconfig:"ICON_EMOJI" default:":robot_face:"`
	ChannelID     string `envconfig:"CHANNEL_ID" default:"#enter_leave_log"`
	CSVLogToken string `envconfig:"CSV_LOG_TOKEN" default:""`
	CSVLogChannelID  string `envconfig:"CSV_LOG_CHANNEL_ID" default:""`
}

var slackinfo SlackInfo

func init() {
	err := envconfig.Process("AGD", &slackinfo)
	if err != nil {
		panic(err)
	}
}

// GetSlackInfo returns slack config info
func GetSlackInfo() SlackInfo {
	return slackinfo
}
