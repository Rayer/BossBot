package BossBot

import (
	"encoding/json"
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
				controller := MsgSchedulerController{conf}
				wm, err = controller.HandleResponse(msgAct)
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
		var ret slack.WebhookMessage
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch s.Command {
		case "/bb_broadcast_list":
			msc := MsgSchedulerController{conf}
			ret, err = msc.HandleRequest()
			if err != nil {
				log.Errorf("Error return : %+v", err)
			}

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		out, err := json.Marshal(ret)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
	})

	log.Println("Starting server....")
	_ = http.ListenAndServe(":5601", nil)
	return nil
}
