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
	Action         string `json:"action"`
	ScheduleItemId int    `json:"item_id"`
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

			//TODO : Make it const
			if msgAct.CallbackId == "MsgSchOperation" {
				//Handler of MsgSchOperation
				//Get action and schedule ID. Usually, there should be only 1 action
				schAction := MsgScheduleIdActions{}
				err := json.Unmarshal([]byte(msgAct.Actions[0].Value), &schAction)
				if err != nil {
					//....should not happen?
					log.Warnln(err)
				}

				mb := MessageBroadcaster{conf}

				w.Header().Set("Content-Type", "application/json")

				switch schAction.Action {
				case "invoke":
					_, err = mb.InvokeBroadcast(schAction.ScheduleItemId)
					if err != nil {
						log.Warnf("Fail at : %+v", schAction)
					}
					//http.Post(msgAct.ResponseUrl, "application/json", strings.NewReader("{\"aaa\":\"invoke22222\"}"))
					_, err = w.Write([]byte("{\"aaa\":\"invoke\"}"))

					break
				case "enable":
					_, err = w.Write([]byte("{\"aaa\":\"enable\"}"))
					break
				case "disable":
					_, err = w.Write([]byte("{\"aaa\":\"disable\"}"))
					break
				}

			}
			return
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
					ScheduleItemId: broadcast.Id,
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
					log.Warnf("Error marshalling json : %+v\n", value)
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
					log.Warnf("Error marshalling json : %+v\n", value)
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
					CallbackID: "MsgSchOperation",
				}

				if params.Attachments == nil {
					params.Attachments = []slack.Attachment{}
				}

				params.Attachments = append(params.Attachments, attachment)

			}

			ret, err := json.Marshal(params)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(ret)

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting server....")
	_ = http.ListenAndServe(":5601", nil)
	return nil
}
