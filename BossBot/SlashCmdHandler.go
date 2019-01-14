package BossBot

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MsgSchedulerController struct {
	config             Configuration
	messageBroadcaster MessageBroadcaster
}

func (msc *MsgSchedulerController) HandleItem(w http.ResponseWriter, msgAct slackevents.MessageAction) (slack.WebhookMessage, error) {

	if msgAct.CallbackId != "MsgSchOperation" {
		return slack.WebhookMessage{Text: "Not correct handler"}, errors.Errorf("Not correct handler")
	}
	//Handler of MsgSchOperation
	//Get action and schedule ID. Usually, there should be only 1 action
	schAction := MsgScheduleIdActions{}
	err := json.Unmarshal([]byte(msgAct.Actions[0].Value), &schAction)
	if err != nil {
		//....should not happen?
		log.Warnln(err)
	}

	mb := msc.messageBroadcaster
	whm := slack.WebhookMessage{}
	w.Header().Set("Content-Type", "application/json")

	switch schAction.Action {
	case "invoke":
		_, err = mb.InvokeBroadcast(schAction.ScheduleItemId)
		if err != nil {
			log.Warnf("Fail at : %+v", schAction)
			whm.Text = fmt.Sprintf("Schedule ID : %d failed to be invoked", schAction.ScheduleItemId)

		} else {
			whm.Text = fmt.Sprintf("Schedule ID : %d is successfully invoked!", schAction.ScheduleItemId)
		}
		break
	case "enable":
		err = mb.SetActive(schAction.ScheduleItemId, true)
		if err != nil {
			log.Warnf("Fail at : %+v", schAction)
			whm.Text = fmt.Sprintf("Schedule ID : %d failed to be enabled", schAction.ScheduleItemId)

		} else {
			whm.Text = fmt.Sprintf("Schedule ID : %d is successfully enabled!", schAction.ScheduleItemId)
		}
		break
	case "disable":
		err = mb.SetActive(schAction.ScheduleItemId, false)
		whm.Text = fmt.Sprintf("Schedule ID : %d successfully disabled!", schAction.ScheduleItemId)
		if err != nil {
			log.Warnf("Fail at : %+v", schAction)
			whm.Text = fmt.Sprintf("Schedule ID : %d failed to be disabled", schAction.ScheduleItemId)

		} else {
			whm.Text = fmt.Sprintf("Schedule ID : %d is successfully disabled!", schAction.ScheduleItemId)
		}
		break
	}
	return whm, nil
}
