package utils

import (
	"github.com/kelseyhightower/envconfig"
)

type SlackInfo struct {
	WEBHOOKURL  string `default:"https://hooks.slack.com/services/"`
	UserName    string `default:"Logging Bot"`
	IconEmoji   string `default:":robot_face:"`
	Channel     string `default:"#enter_leave_log"`
	CSVLOGTOKEN string `default:""`
	CSVLOGCHID  string `default:""`
}

func (sl SlackInfo) GetWebHookURL() string {
	return sl.WEBHOOKURL
}

var slackinfo SlackInfo

func init() {
	slackinfo, _ = Init()
}

func Init() (SlackInfo, error) {
	err := envconfig.Process("AGD", &slackinfo)
	if err != nil {
		return SlackInfo{}, err
	}
	return slackinfo, nil
}

func GetSlackInfo() SlackInfo {
	return slackinfo
}
