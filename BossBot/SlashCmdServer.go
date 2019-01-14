package BossBot

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
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
		w.Header().Set("Content-Type", "application/json")
		var wm slack.WebhookMessage
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
				controller := MsgSchedulerController{conf, MessageBroadcaster{conf}}
				wm, err = controller.HandleItem(w, msgAct)
				if err != nil {
					log.Errorf("Error handling MsgSchOperation, error : %s", err)
					wm.Text = "Error handling MsgSchOperation!"
				}
			}

			ret, _ := json.Marshal(wm)
			w.Write(ret)
			return
		}

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

			log.Debugf("Incoming cmd bb_broadcast_list")

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

				var fields []slack.AttachmentField

				//Makeup id / message_id / Message field
				fields = append(fields, slack.AttachmentField{
					Title: "Scheduler ID",
					Value: strconv.Itoa(broadcast.Id),
					Short: true,
				})

				fields = append(fields, slack.AttachmentField{
					Title: "Message ID",
					Value: strconv.Itoa(broadcast.MessageId),
					Short: true,
				})

				fields = append(fields, slack.AttachmentField{
					Title: "Message",
					Value: broadcast.Message,
					Short: false,
				})

				//Start date and End date

				fields = append(fields, slack.AttachmentField{
					Title: "Start From",
					Value: func() string {
						if broadcast.StartDate.Valid {
							return broadcast.StartDate.Time.Format(time.RFC822)
						}
						return "Not set"
					}(),
					Short: true,
				})

				fields = append(fields, slack.AttachmentField{
					Title: "Ends at",
					Value: func() string {
						if broadcast.EndDate.Valid {
							return broadcast.StartDate.Time.Format(time.RFC822)
						}
						return "Not set"
					}(),
					Short: true,
				})

				//Recursive date time

				fields = append(fields, slack.AttachmentField{
					Title: "Run at nth day of month",
					Value: func() string {
						if broadcast.Day.Valid {
							return "Every " + strconv.Itoa(int(broadcast.Day.Int64)) + " of the month"
						}
						return "Not set"
					}(),
					Short: true,
				})

				fields = append(fields, slack.AttachmentField{
					Title: "Run at day in week",
					Value: func() string {
						if broadcast.WeekDay.Valid {
							return "Every " + time.Weekday(broadcast.WeekDay.Int64).String()
						}
						return "Not set"
					}(),
					Short: true,
				})

				fields = append(fields, slack.AttachmentField{
					Title: "Broadcast Time",
					Value: broadcast.BroadcastTime,
					Short: true,
				})

				fields = append(fields, slack.AttachmentField{
					Title: "Channel",
					Value: broadcast.ChannelName,
					Short: true,
				})

				attachment := slack.Attachment{
					Text:       "Message detail",
					Actions:    actions,
					Color:      color,
					CallbackID: "MsgSchOperation",
					Fields:     fields,
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
