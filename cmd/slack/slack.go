package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
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

// SlackNotify sends slack notification.
func SlackNotify(Name string, UID string, isEnter bool, Timestamp time.Time, Ext string) {
	HookJSON := &WebHook{
		Username:  "Logging Bot",
		IconEmoji: ":robot_face:",
		Channel:   "#enter_leave_log",
	}

	if isEnter {
		HookJSON.Text = fmt.Sprintf("%v : %v さんが %v に入室しました。", UID, Name, Timestamp)
	} else {
		HookJSON.Text = fmt.Sprintf("%v : %v さんが %v に退室しました。", UID, Name, Timestamp)

		var RawJSON = []byte(Ext)
		var ExtList = make(map[string]interface{})

		err := json.Unmarshal(RawJSON, &ExtList)
		if err != nil {
			log.Println(err)
			log.Println("Error has occured in SlackNotify. User Ext is : ", Ext)
			errOccuredSlackNotify()
			return
		}

		useage := ExtList["Use"].([]interface{})
		mess := ExtList["message"].(string)

		At := Attachment{
			Title: "アンケート結果",
			Text:  fmt.Sprintf("目的 : %v \n 感想 : %v", useage, mess),
		}

		HookJSON.Attachments = append(HookJSON.Attachments, At)
	}
	err := postEnterLeaveLog(HookJSON)
	if err != nil {
		log.Println(err)
		errOccuredSlackNotify()
		return
	}
}

func errOccuredSlackNotify() {
	log.Println("error")
}

func postEnterLeaveLog(ellog *WebHook) error {

	IncomingURL := os.Getenv("SLACK_WEBHOOK_URL")

	params, _ := json.Marshal(ellog)

	resp, _ := http.PostForm(
		IncomingURL,
		url.Values{"payload": {string(params)}},
	)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	println(string(body))
	return nil
}
