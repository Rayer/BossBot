package ChatBot

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

type RootScenario struct {
	DefaultScenarioImpl
}

func (rs *RootScenario) InitScenario(uc *UserContext) error {
	rs.ScenarioSelf = rs
	rs.registerState("entry", &EntryState{}, rs)
	rs.registerState("second", &SecondState{}, rs)
	return nil
}

func (rs *RootScenario) EnterScenario(source Scenario) error {
	log.Debugln("Entering root scenario")
	return nil
}

func (rs *RootScenario) ExitScenario(askFrom Scenario) error {
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
	DefaultScenarioStateImpl
}

func (es *EntryState) RenderMessage() (string, error) {
	return "Hey it's BossBot! Are you going to [submit report], [manage broadcasts] or [check]?", nil
}

func (es *EntryState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "submit report") {
		es.parent.changeStateByName("second")
		return "Exit with 1", nil
	} else if strings.Contains(input, "manage broadcast") {
		es.parent.changeStateByName("second")
		return "Exit with 2", nil
	}

	return "Nothing done", nil
}

func (es *EntryState) GetParentScenario() Scenario {
	return es.parent
}

type SecondState struct {
	DefaultScenarioStateImpl
}

func (ss *SecondState) RenderMessage() (string, error) {
	return "This is second message, you can only [exit] in order to get out of here", nil
}

func (ss *SecondState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "exit") {
		ss.parent.changeStateByName("entry")
		return "Exiting...", nil
	}
	return "Not exit, stay here.", nil
}

func (ss *SecondState) GetParentScenario() Scenario {
	return ss.parent
}

func (rs *RootScenario) Name() string {
	return "RootScenario"
}

func (rs *RootScenario) GetUserContext() *UserContext {
	return nil
}
