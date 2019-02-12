package BossBot

import (
	"ChatBot"
	"bytes"
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

	slack_client := conf.ServiceContext.SlackClient
	//slack_rtm := conf.ServiceContext.SlackRTM
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

	//Events API
	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Debugf("Get event message : %s", buf)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{conf.SlackVerifyToken}))
		log.Debugf("Parse into events api event : %+v", eventsAPIEvent)
		if e != nil {
			log.Debugf("Error occured : %s", e)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			postParams := slack.PostMessageParameters{}
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				slack_client.PostMessage(ev.Channel, "Yes, hello.", postParams)
			case *slackevents.MessageAction:
				//Need to render a message while open dm

			case *slackevents.MessageEvent:
				//Do bot part
				msgevent := innerEvent.Data.(*slackevents.MessageEvent)

				//Discard message from bot
				if msgevent.User == "" && msgevent.BotID != "" {
					return
				}

				log.Debugf("Got message from user : %s  botid : %s with message : %s", msgevent.User, msgevent.BotID, msgevent.Text)

				//Translate it to name
				var name string
				userProfile, err := slack_client.GetUserProfile(msgevent.User, false)
				if err != nil {
					log.Warnf("Can't translate slack uid %s to name, use uid instead of name", msgevent.User)
					log.Warnf("error : %+v", err)
					name = msgevent.User
				} else {
					log.Debugf("Found %s as %s", name, msgevent.User)
					name = userProfile.DisplayName
				}

				chatBot := conf.ServiceContext.ChatBotClient
				userContext := chatBot.GetUserContext(name)
				if userContext == nil {
					userContext = chatBot.CreateUserContext(name, func() ChatBot.Scenario {
						return &RootScenario{}
					})
				}

				handledMessage, _ := userContext.HandleMessage(msgevent.Text)
				if handledMessage != "" {
					slack_client.PostMessage(msgevent.Channel, handledMessage, postParams)
				}

				ret, _ := userContext.RenderMessage()

				slack_client.PostMessage(msgevent.Channel, ret, postParams)

				return

			}
		}
	})

	log.Println("Starting server....")
	_ = http.ListenAndServe(":5601", nil)
	return nil
}
