package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebHook is structure for slack notify
type WebHook struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment is structure for slack notify data
type Attachment struct {
	Title   string `json:"title"`
	Pretext string `json:"pretext"`
	Color   string `json:"color"`
	Text    string `json:"text"`
}

// WebHookInit returns new WebHook data
func WebHookInit(cfg *SlackInfo) *WebHook {
	return &WebHook{
		Username:  cfg.UserName,
		IconEmoji: cfg.IconEmoji,
		Channel:   cfg.Channel,
	}
}

// SlackNotify sends slack notification.
func SlackNotify(Name string, UID string, isEnter bool, Timestamp time.Time, Ext string) error {

	cfg := GetSlackInfo()
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

		var purpose string
		for _, w := range useage {
			purpose += fmt.Sprintf("%v, ", w)
		}
		purpose = strings.TrimSuffix(purpose, ", ")

		HookJSON.Attachments = append(HookJSON.Attachments, Attachment{
			Title: "アンケート結果",
			Text:  fmt.Sprintf("目的 : %v \n 感想 : %v", purpose, mess),
			Color: "#FF9500",
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
	IncomingURL := GetSlackInfo().GetWebHookURL()
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

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
