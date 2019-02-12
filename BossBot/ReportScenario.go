package BossBot

import (
	"ChatBot"
	"Utilities"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ReportScenario struct {
	ChatBot.DefaultScenarioImpl
	ThisWeekInDev []string
	ThisWeekDone  []string
}

func (rs *ReportScenario) InitScenario(uc *ChatBot.UserContext) error {
	rs.DefaultScenarioImpl.InitScenario(uc)
	rs.RegisterState("entry", &ReportEntryState{}, rs)
	rs.RegisterState("creating_done", &ReportCreatingDone{}, rs)
	rs.RegisterState("creating_indev", &ReportCreatingInDev{}, rs)
	rs.RegisterState("confirm", &ReportConfirm{}, rs)
	return nil
}

func (rs *ReportScenario) EnterScenario(source ChatBot.Scenario) error {
	return nil
}

func (rs *ReportScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	return nil
}

func (rs *ReportScenario) DisposeScenario() error {
	return nil
}

func (rs *ReportScenario) Name() string {
	return "Weekly Report Scenario"
}

/*
States :
1. Entry - Greeting with current period report or re-create, if not, [Create Report]
2. CreatingDone
3. CreatingInDev
4. Review
*/

type ReportEntryState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (res *ReportEntryState) RenderMessage() (string, error) {
	/*
		Designed functionality :
		1. Let user view logs before (Not in this prototype)
		2. Show if log is submitted in this week. If so, show it and ask if it need to be recreate or exit
		3. If no report in this week, ask user to create one
	*/

	conf := GetConfiguration()
	db := conf.ServiceContext.DBObject.GetDB()
	slack_id := res.GetParentScenario().GetUserContext().User
	year, week := time.Now().ISOWeek()
	result, err := db.Query("select * from bb_weekly_report where user_slack_id = ? and year = ? and week_of_year = ?", slack_id, year, week)
	if err != nil {
		return "", errors.Wrap(err, "Error executing RenderMessage")
	}

	var retText string
	if result.Next() {
		var weeklyReportItem WeeklyReportItem
		err = Utilities.RowsToStruct("bb_data", result, &weeklyReportItem)
		if err != nil {
			return "", errors.Wrap(err, "Fail to marshal datatype : WeeklyReportItem")
		}

		retText = fmt.Sprintf("Year %d week %d report --- \n Done : \n %s \n On Going : \n %s \n Would you like to [create report] or [exit]?", weeklyReportItem.Year, weeklyReportItem.WeekOfYear, weeklyReportItem.Done, weeklyReportItem.OnGoing)

	} else {
		retText = "Hey, we don't see logs this week. Would you like to [create report]? or [view reports] in previous weeks? You also can [exit] if no longer need to operating with logs"
	}

	return retText, nil
}

func (res *ReportEntryState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "create report") {
		res.GetParentScenario().ChangeStateByName("creating_done")
		return "Ok let's creating a report", nil
	} else if strings.Contains(input, "view report") {
		return "Not really implemented in this prototype version... maybe later", nil
	} else if strings.Contains(input, "exit") {
		res.GetParentScenario().GetUserContext().ReturnLastScenario()
		return "Let's back to previous session", nil
	}

	return "I don't really understand.... can you use another phrase with same meaning?", nil
}

type ReportCreatingDone struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rcd *ReportCreatingDone) RenderMessage() (string, error) {
	return "What task have been done in this week? or there is [good for now]?", nil
}

func (rcd *ReportCreatingDone) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "good for now") {
		rcd.GetParentScenario().ChangeStateByName("creating_indev")
		return "Done in done", nil
	}

	doneList := rcd.GetParentScenario().(*ReportScenario).ThisWeekDone
	rcd.GetParentScenario().(*ReportScenario).ThisWeekDone = append(doneList, input)

	return "Recorded (done) : " + input, nil
}

type ReportCreatingInDev struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rcid *ReportCreatingInDev) RenderMessage() (string, error) {
	return "What task is in dev this week? or it's [good for now]?", nil
}

func (rcid *ReportCreatingInDev) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "good for now") {
		rcid.GetParentScenario().ChangeStateByName("confirm")
		return "Done in dev", nil
	}

	indevList := rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev
	rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev = append(indevList, input)

	return "Recorded (On Going): " + input, nil
}

type ReportConfirm struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rc *ReportConfirm) RenderMessage() (string, error) {
	doneList := rc.GetParentScenario().(*ReportScenario).ThisWeekDone
	indevList := rc.GetParentScenario().(*ReportScenario).ThisWeekInDev

	ret := "Will you [submit] or [discard] follow report entries : "
	ret += "Done : \n"
	for _, done := range doneList {
		ret += done + "\n"
	}

	ret += "In Dev : \n"
	for _, inDev := range indevList {
		ret += inDev + "\n"
	}

	return ret, nil

}

func (rc *ReportConfirm) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "submit") {
		err := rc.submitResult()
		if err != nil {
			log.Errorf("Error : %+v", err)
			return "Error submitting report!", errors.Wrap(err, "Error submitting report!")
		}
		rc.GetParentScenario().GetUserContext().ReturnLastScenario()
		return "Submitted", nil
	} else if strings.Contains(input, "discard") {
		rc.GetParentScenario().GetUserContext().ReturnLastScenario()
		return "Discarded", nil
	}

	return "I don't really understand.....", nil
}

func (rc *ReportConfirm) submitResult() error {
	parent := rc.GetParentScenario().(*ReportScenario)
	var done string
	var ongoing string

	for _, d := range parent.ThisWeekDone {
		done += d
		done += "\n"
	}

	for _, o := range parent.ThisWeekInDev {
		ongoing += o
		ongoing += "\n"
	}

	slack_id := parent.GetUserContext().User
	db := GetConfiguration().ServiceContext.DBObject.GetDB()
	year, week := time.Now().ISOWeek()

	_, err := db.Exec("insert into bb_weekly_report (year, week_of_year, user_slack_id, done, ongoing) values (?, ?, ?, ?, ?) ", year, week, slack_id, done, ongoing)
	if err != nil {
		return errors.Wrap(err, "Error submitting result!")
	}

	return nil
}
