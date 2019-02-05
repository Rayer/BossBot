package main

import (
	"ChatBot"
	"strings"
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

	return "Hey, we don't see logs this week. Would you like to [create report]? or [view reports] in previous weeks? You also can [exit] if no longer need to operating with logs", nil
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

	return "Recorded (indev): " + input, nil
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
		rc.GetParentScenario().GetUserContext().ReturnLastScenario()
		return "Submitted", nil
	} else if strings.Contains(input, "discard") {
		rc.GetParentScenario().GetUserContext().ReturnLastScenario()
		return "Discarded", nil
	}

	return "I don't really understand.....", nil
}
