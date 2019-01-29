package ChatBot

import "strings"

type RootScenario struct{
	DefaultScenarioImpl
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {

}

func (es *EntryState) RenderMessage() string {
	return "This is an demo scenario and state. Are you going to invoke [first] one scenario or [second] scenario?"
}

func (es *EntryState) HandleMessage(input string) string {
	if strings.Contains(input, "first") {

	} else
	if strings.Contains(input, "second") {

	}

	return "Nothing done"
}

func (es *EntryState) GetParentScenario() *Scenario {
	panic("implement me")
}


type SecondState struct {

}

func (ss *SecondState) RenderMessage() string {
	panic("implement me")
}

func (ss *SecondState) HandleMessage(input string) string {
	panic("implement me")
}

func (ss *SecondState) GetParentScenario() *Scenario {
	panic("implement me")
}

func (rs *RootScenario) Name() string {
	return "RootScenario"
}
