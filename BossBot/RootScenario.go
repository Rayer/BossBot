package BossBot

import (
	"ChatBot"
	. "SlackChatBot"
	log "github.com/sirupsen/logrus"
	"strings"
)

type RootScenario struct {
	ChatBot.DefaultScenarioImpl
	SlackScenarioImpl
}

func (rs *RootScenario) InitScenario(uc *ChatBot.UserContext) error {
	rs.DefaultScenarioImpl.InitScenario(uc)
	rs.SlackScenarioImpl.InitSlackScenario(rs)
	rs.RegisterState("entry", &EntryState{}, rs)
	rs.RegisterState("second", &SecondState{}, rs)
	return nil
}

func (rs *RootScenario) EnterScenario(source ChatBot.Scenario) error {
	log.Debugln("Entering root scenario")
	return nil
}

func (rs *RootScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	log.Debugln("Exiting root scenario")
	return nil
}

func (rs *RootScenario) DisposeScenario() error {
	log.Debugln("Disposing root scenario")
	return nil
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {
	ChatBot.DefaultScenarioStateImpl
	SlackScenarioStateImpl
	name string
}

func (es *EntryState) InitScenarioState(scenario ChatBot.Scenario) {
	es.name = "EntryState"
	es.SlackScenarioStateImpl = *NewSlackScenarioStateImpl(es)
	es.KeywordHandler().RegisterKeyword(&Keyword{
		Keyword: "manage weekly reports",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.GetUserContext().InvokeNextScenario(&ReportScenario{}, ChatBot.Stack)
			return "Go to report scenario", nil
		},
	})

	es.KeywordHandler().RegisterKeyword(&Keyword{
		Keyword: "manage broadcasts",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.ChangeStateByName("second")
			return "Under construction!", nil
		},
	})

	es.KeywordHandler().RegisterKeyword(&Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			return "Hey it is BossBot! How can I serve you?", nil
		},
	})
}

func (es *EntryState) RenderMessage() (string, error) {
	return "Are you going to [manage weekly reports], [manage broadcasts] or [check]?", nil
}

func (es *EntryState) HandleMessage(input string) (string, error) {

	ret, err := es.KeywordHandler().ParseAction(input)
	if err != nil {
		return "Error handling message!", err
	}
	return ret, nil
}

type SecondState struct {
	ChatBot.DefaultScenarioStateImpl
	SlackScenarioStateImpl
}

func (ss *SecondState) InitScenarioState(scenario ChatBot.Scenario) {
	ss.SlackScenarioStateImpl = *NewSlackScenarioStateImpl(ss)
	ss.KeywordHandler().RegisterKeyword(&Keyword{
		Keyword: "exit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.ChangeStateByName("entry")
			return "Exiting...", nil
		},
	})
}

func (ss *SecondState) RenderMessage() (string, error) {
	return "This page is under construction, you can [exit] to last scene", nil
}

func (ss *SecondState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "exit") {
		ss.GetParentScenario().ChangeStateByName("entry")
		return "Exiting...", nil
	}
	return "Not exit, stay here.", nil
}

func (rs *RootScenario) Name() string {
	return "RootScenario"
}
