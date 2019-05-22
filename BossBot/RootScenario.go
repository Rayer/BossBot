package BossBot

import (
	"github.com/Rayer/chatbot"
	log "github.com/sirupsen/logrus"
)

type RootScenario struct {
	ChatBot.DefaultScenarioImpl
}

func (rs *RootScenario) InitScenario(uc *ChatBot.UserContext) error {
	rs.DefaultScenarioImpl.InitScenario(uc)
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
	name string
}

func (es *EntryState) InitScenarioState(scenario ChatBot.Scenario) {
	es.Init(scenario, es)
	es.name = "EntryState"
	es.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "manage weekly reports",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.GetUserContext().InvokeNextScenario(&ReportScenario{}, ChatBot.Stack)
			return "Go to report scenario", nil
		},
	})

	es.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "manage broadcasts",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.ChangeStateByName("second")
			return "Under construction!", nil
		},
	})

	es.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			return "Hey it is BossBot! How can I serve you?", nil
		},
	})
}

func (es *EntryState) RawMessage() (string, error) {
	return "Are you going to [manage weekly reports], [manage broadcasts] or [check]?", nil
}

type SecondState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (ss *SecondState) InitScenarioState(scenario ChatBot.Scenario) {
	ss.Init(scenario, ss)
	ss.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "exit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
			scenario.ChangeStateByName("entry")
			return "Exiting...", nil
		},
	})
}

func (ss *SecondState) RawMessage() (string, error) {
	return "This page is under construction, you can [exit] to last scene", nil
}


func (rs *RootScenario) Name() string {
	return "RootScenario"
}
