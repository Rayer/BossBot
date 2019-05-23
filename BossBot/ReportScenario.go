package BossBot

import (
	"Utilities"
	"fmt"
	"github.com/Rayer/chatbot"
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
	rs.RegisterState("gather_report", &GatherReport{}, rs)
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

func (res *ReportEntryState) InitScenarioState(scenario ChatBot.Scenario) {
	res.Init(scenario, res)
	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "create report",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			_ = res.GetParentScenario().ChangeStateByName("creating_done")
			return "Ok let's creating a report", nil
		},
	})

	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "view report",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			return "Not really implemented in this prototype version... maybe later", nil
		},
	})

	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "exit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			_ = res.GetParentScenario().GetUserContext().ReturnLastScenario()
			return "Let's back to previous session", nil
		},
	})

	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "Gather this week reports",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			res.GetParentScenario().ChangeStateByName("gather_report")
			return "Let's gather reports", nil
		},
	})
}

func (res *ReportEntryState) RawMessage() (string, error) {
	/*
		Designed functionality :
		1. Let user view logs before (Not in this prototype)
		2. Show if log is submitted in this week. If so, show it and ask if it need to be recreate or exit
		3. If no report in this week, ask user to create one
	*/

	conf := GetConfiguration()
	db := conf.ServiceContext.DBObject.GetDB()
	name := res.GetParentScenario().GetUserContext().User
	year, week := time.Now().ISOWeek()
	result, err := db.Query("select * from bb_weekly_report where user_id = ? and year = ? and week_of_year = ?", name, year, week)
	if err != nil {
		return "", errors.Wrap(err, "Error executing RenderMessage")
	}

	var retText string
	defer func() {
		err := result.Close()
		if err != nil {
			log.Warnf("Error closing query row : %+v", err)
		}
	}()
	if result.Next() {
		var weeklyReportItem WeeklyReportItem
		err = Utilities.RowsToStruct("bb_data", result, &weeklyReportItem)
		if err != nil {
			return "", errors.Wrap(err, "Fail to marshal datatype : WeeklyReportItem")
		}

		retText = fmt.Sprintf("Hey %s, I see you have report : \nYear %d week %d report --- \nDone : \n%s\nOn Going :\n%s\n Would you like to [create report], [Gather this week reports] or [exit]?", name, weeklyReportItem.Year, weeklyReportItem.WeekOfYear, weeklyReportItem.Done, weeklyReportItem.OnGoing)

	} else {
		retText = fmt.Sprintf("Hey %s, we don't see logs this week. Would you like to [create report]? [Gather this week reports] or [view reports] in previous weeks? You also can [exit] if no longer need to operating with logs", name)
	}

	return retText, nil
}

type ReportCreatingDone struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rcd *ReportCreatingDone) InitScenarioState(scenario ChatBot.Scenario) {
		rcd.Init(scenario, rcd)
		rcd.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "good for now",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			_ = rcd.GetParentScenario().ChangeStateByName("creating_indev")
			doneTasks := "*Tasks done for this week* : \n"
			for _, task := range rcd.GetParentScenario().(*ReportScenario).ThisWeekDone {
				doneTasks += " - " + task
			}

			return doneTasks, nil
		},
	})

	//Register default action.
	rcd.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			doneList := rcd.GetParentScenario().(*ReportScenario).ThisWeekDone
			splited := strings.Split(input, "\n")
			for _, task := range splited {
				doneList = append(doneList, task+"\n")
			}

			rcd.GetParentScenario().(*ReportScenario).ThisWeekDone = doneList
			return "Recorded (done) : \n" + input, nil
		},
	})
}

func (rcd *ReportCreatingDone) RawMessage() (string, error) {
	return "What task *have been done* in this week? or there is [good for now]?", nil
}


type ReportCreatingInDev struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rcid *ReportCreatingInDev) InitScenarioState(scenario ChatBot.Scenario) {
	rcid.Init(scenario, rcid)
	rcid.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "good for now",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			rcid.GetParentScenario().ChangeStateByName("confirm")

			inprogressTasks := "*Tasks in progress for this week* : \n"
			for _, task := range rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev {
				inprogressTasks += " - " + task
			}

			return inprogressTasks, nil
			//return "Done in dev", nil
		},
	})

	rcid.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			indevList := rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev
			splited := strings.Split(input, "\n")
			for _, task := range splited {
				indevList = append(indevList, task+"\n")
			}
			rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev = indevList
			return "Recorded (In Progress): \n" + input, nil
		},
	})
}

func (rcid *ReportCreatingInDev) RawMessage() (string, error) {
	return "What task is *in progress* this week? or it's [good for now]?", nil
}


type ReportConfirm struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rc *ReportConfirm) InitScenarioState(scenario ChatBot.Scenario) {
	rc.Init(scenario, rc)
	rc.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "submit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			err := rc.submitResult()
			if err != nil {
				log.Errorf("Error : %+v", err)
				return "Error submitting report!", errors.Wrap(err, "Error submitting report!")
			}
			_ = rc.GetParentScenario().GetUserContext().ReturnLastScenario()
			return "Submitted", nil
		},
	})
	rc.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "discard",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			_ = rc.GetParentScenario().GetUserContext().ReturnLastScenario()
			return "Discarded", nil
		},
	})
}

func (rc *ReportConfirm) RawMessage() (string, error) {
	doneList := rc.GetParentScenario().(*ReportScenario).ThisWeekDone
	indevList := rc.GetParentScenario().(*ReportScenario).ThisWeekInDev

	ret := "Will you [submit] or [discard] follow report entries : \n"
	ret += "*Done* : \n"
	for _, done := range doneList {
		ret += done + "\n"
	}

	ret += "*In Progress* : \n"
	for _, inDev := range indevList {
		ret += inDev + "\n"
	}

	return ret, nil

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

	user := parent.GetUserContext().User
	db := GetConfiguration().ServiceContext.DBObject.GetDB()
	year, week := time.Now().ISOWeek()

	//Delete old entry
	_, err := db.Exec("delete from bb_weekly_report where year = ? and week_of_year = ? and user_id = ?", year, week, user)
	if err != nil {
		return errors.Wrap(err, "Error deleting old entry!")
	}
	_, err = db.Exec("insert into bb_weekly_report (year, week_of_year, user_id, done, ongoing) values (?, ?, ?, ?, ?) ", year, week, user, done, ongoing)
	if err != nil {
		return errors.Wrap(err, "Error submitting result!")
	}

	return nil
}

type GatherReport struct {
	ChatBot.DefaultScenarioStateImpl
}

func (gr *GatherReport) InitScenarioState(scenario ChatBot.Scenario) {
	gr.Init(scenario, gr)
	gr.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "by person",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			year, week := time.Now().ISOWeek()
			reports, err := GetWeeklyReports(year, week)
			if err != nil {
				return "Error getting reports!", err
			}

			var output string

			for _, report := range reports {
				output += fmt.Sprintf("%s :\n Completed :\n%s\n In Progress : \n%s\n\n------------\n", report.UserSlackId, report.Done, report.OnGoing)
			}

			return output, nil
		},
	})

	gr.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "by progress",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			year, week := time.Now().ISOWeek()
			reports, err := GetWeeklyReports(year, week)
			if err != nil {
				return "Error getting reports!", err
			}

			var completed string
			var indev string
			var submittedUser string

			for _, report := range reports {
				completed += report.Done
				indev += report.OnGoing
				submittedUser += report.UserSlackId + " "
			}

			return fmt.Sprintf("Year : %d Week : %d Weekly Report : \n\nCompleted:\n%s\nIn Progress:\n%s\nSubmitted user : %s\n", year, week, completed, indev, submittedUser), nil
		},
	})

	gr.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "exit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			gr.GetParentScenario().ChangeStateByName("entry")
			return "Return to last scene", nil
		},
	})

}

func (gr *GatherReport) RawMessage() (string, error) {
	return "Would you gather reports [by person], [by progress] or [exit]?", nil
}
