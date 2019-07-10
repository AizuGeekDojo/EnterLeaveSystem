package utils

import (
	"github.com/kelseyhightower/envconfig"
)

// SlackInfo is structure for slack info in env
type SlackInfo struct {
	WEBHOOKURL  string `default:"https://hooks.slack.com/services/"`
	UserName    string `default:"Logging Bot"`
	IconEmoji   string `default:":robot_face:"`
	Channel     string `default:"#enter_leave_log"`
	CSVLOGTOKEN string `default:""`
	CSVLOGCHID  string `default:""`
}

// GetWebHookURL returns slack webhook URL
func (sl SlackInfo) GetWebHookURL() string {
	return sl.WEBHOOKURL
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
