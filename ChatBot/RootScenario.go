package ChatBot

import "strings"

type RootScenario struct {
	DefaultScenarioImpl
}

func (rs *RootScenario) InitScenario(uc *UserContext) error {

	es := EntryState{}
	ss := SecondState{}
	es.InitWithParent(rs)
	ss.InitWithParent(rs)
	rs.registerState("entry", &es)
	rs.registerState("second", &ss)
	return nil
}

func (rs *RootScenario) EnterScenario(source Scenario) error {
	//Init states
	return nil
}

func (rs *RootScenario) ExitScenario(askFrom Scenario) error {
	return nil
}

func (rs *RootScenario) DisposeScenario() error {
	panic("implement me")
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {
	parent Scenario
}

func (es *EntryState) InitWithParent(parent Scenario) error {
	es.parent = parent
	es.parent.registerState("entry", es)
	return nil
}

func (es *EntryState) RenderMessage() (string, error) {
	return "This is an demo scenario and state. Are you going to invoke [first] one scenario or [second] scenario?", nil
}

func (es *EntryState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "first") {
		es.parent.changeStateByName("second")
		return "Exit with 1", nil
	} else if strings.Contains(input, "second") {
		es.parent.changeStateByName("second")
		return "Exit with 2", nil
	}

	return "Nothing done", nil
}

func (es *EntryState) GetParentScenario() Scenario {
	return es.parent
}

type SecondState struct {
	parent Scenario
}

func (ss *SecondState) InitWithParent(parent Scenario) error {
	ss.parent = parent
	parent.registerState("second", ss)
	return nil
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
