package BossBot

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type MsgSchedulerController struct {
	config Configuration
	//messageBroadcaster MessageBroadcaster
}

func (msc *MsgSchedulerController) HandleResponse(msgAct slackevents.MessageAction) (slack.WebhookMessage, error) {

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

	mb := MessageBroadcaster{msc.config}
	whm := slack.WebhookMessage{}

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

func (msc *MsgSchedulerController) HandleRequest() (slack.WebhookMessage, error) {

	log.Debugf("Incoming cmd bb_broadcast_list")

	errMsg := slack.WebhookMessage{}

	bdList, err := GetBroadcastList()

	if err != nil {
		log.Warnf("%s : %s", err, "Error getting Broadcast List")
		errMsg.Text = "Error getting Broadcast List!"
		return errMsg, err
	}

	ret := slack.WebhookMessage{Text: "Here is broadcast items in system :"}

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
			errMsg.Text = "Error marshalling json"
			return errMsg, err
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
			errMsg.Text = "Error marshalling json"
			return errMsg, err
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

		if ret.Attachments == nil {
			ret.Attachments = []slack.Attachment{}
		}

		ret.Attachments = append(ret.Attachments, attachment)

	}

	return ret, nil

	//ret, err := json.Marshal(params)
	//w.Header().Set("Content-Type", "application/json")
	//_, _ = w.Write(ret)
}
