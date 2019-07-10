package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/config"
)

type WebHook struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Title   string `json:"title"`
	Pretext string `json:"pretext"`
	Text    string `json:"text"`
}

func WebHookInit(cfg *config.SlackInfo) *WebHook {
	return &WebHook{
		Username:  cfg.UserName,
		IconEmoji: cfg.IconEmoji,
		Channel:   cfg.Channel,
	}
}

// Notify sends slack notification.
func Notify(Name string, UID string, isEnter bool, Timestamp time.Time, Ext string) error {

	cfg := config.GetSlackInfo()
	HookJSON := WebHookInit(&cfg)

	var io string
	if isEnter {
		io = "入室"
	} else {
		io = "退室"
	}

	HookJSON.Text = fmt.Sprintf("%v : %v さんが %v に%vしました。", UID, Name, Timestamp.Format("2006-01-02 15:04:05"), io)

	if Ext != "" {
		var RawJSON = []byte(Ext)
		var ExtList = make(map[string]interface{})

		err := json.Unmarshal(RawJSON, &ExtList)
		if err != nil {
			log.Println(err)
			log.Println("Error has occured in SlackNotify. User Ext is : ", Ext)
			return err
		}

		useage := ExtList["Use"].([]interface{})
		mess := ExtList["message"].(string)

		HookJSON.Attachments = append(HookJSON.Attachments, Attachment{
			Title: "アンケート結果",
			Text:  fmt.Sprintf("目的 : %v \n 感想 : %v", useage, mess),
		})
	}

	err := postEnterLeaveLog(HookJSON)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func postEnterLeaveLog(ellog *WebHook) error {
	IncomingURL := config.GetSlackInfo().GetWebHookURL()
	if IncomingURL == "" {
		return errors.New("Slack URL is not defined")
	}

	params, err := json.Marshal(ellog)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	resp, err := client.PostForm(
		IncomingURL,
		url.Values{"payload": {string(params)}},
	)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	println(string(body))
	return nil
}
