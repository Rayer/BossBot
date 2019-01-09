package BossBot

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MsgScheduleIdActions struct {
	Action string `json:"action"`
	MsgId  int    `json:"msg_id"`
}

func RespServer(conf Configuration) error {

	http.HandleFunc("/slack/interactive", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, item := range r.PostForm["payload"] {
			msgAct := slackevents.MessageAction{}
			err := json.Unmarshal([]byte(item), &msgAct)
			if err != nil {
				log.Errorln("Error handling income interactive message : ", item)
				log.Errorln("Error is : ", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte("{\"aaa\":\"ccc\"}"))
		if err != nil {
			log.Errorln("Fail to send message : ", err)
		}

	})

	http.HandleFunc("/slack/slash_cmds", func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Incoming : %+v \n", r)
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch s.Command {
		case "/bb_broadcast_list":

			//params := &slack.Msg{Text: s.Text}

			//b, err := json.Marshal(params)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			dm, err := GetDataManager(&conf)
			if err != nil {
				log.Errorf("%s : %s", err, "Error getting data manager!")
				return
			}

			bdList, err := dm.GetBroadcastList()

			if err != nil {
				log.Warnf("%s : %s", err, "Error getting Broadcast List")
				return
			}

			log.Tracef("List : %v", bdList)

			params := slack.WebhookMessage{Text: "Here is broadcast items in system :"}

			for _, broadcast := range bdList {

				var actions []slack.AttachmentAction

				var actionBtnName string
				var color string
				value := MsgScheduleIdActions{
					MsgId: broadcast.MessageId,
				}

				if broadcast.Active == 1 {
					actionBtnName = "Disable"
					color = "#3AA3E3"
					value.Action = "disable"
				} else {
					actionBtnName = "Enable"
					color = "#FF0000"
					value.Action = "enable"
				}

				valueJson, err := json.Marshal(&value)

				if err != nil {
					log.Warnln("Error marshelling json : %+v", value)
					return
				}

				actions = append(actions, slack.AttachmentAction{
					Name:  actionBtnName,
					Text:  actionBtnName,
					Type:  "button",
					Style: "primary",
					Value: string(valueJson),
				})

				value.Action = "invoke"
				valueJson, err = json.Marshal(&value)

				if err != nil {
					log.Warnln("Error marshelling json : %+v", value)
					return
				}
				actions = append(actions, slack.AttachmentAction{
					Name:  "Invoke Now",
					Text:  "Invoke Now",
					Type:  "button",
					Style: "danger",
					Value: string(valueJson),
				})

				attachment := slack.Attachment{
					Text:       fmt.Sprintf("%s", broadcast.StringForSlackItem()),
					Actions:    actions,
					Color:      color,
					CallbackID: "MsgSchModify",
				}

				if params.Attachments == nil {
					params.Attachments = []slack.Attachment{}
				}

				params.Attachments = append(params.Attachments, attachment)

			}

			ret, err := json.Marshal(params)
			w.Header().Set("Content-Type", "application/json")
			w.Write(ret)

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting server....")
	http.ListenAndServe(":5601", nil)
	return nil
}
