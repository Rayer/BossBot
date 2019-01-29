package ChatBot

import "strings"

type RootScenario struct{
	DefaultScenarioImpl
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {
	parent *Scenario
}

func (es *EntryState) InitWithParent(parent *Scenario) error {
	es.parent = parent
	(*es.parent).registerState("entry", es)
	return nil
}

func (es *EntryState) RenderMessage() string {
	return "This is an demo scenario and state. Are you going to invoke [first] one scenario or [second] scenario?"
}

func (es *EntryState) HandleMessage(input string) string {
	if strings.Contains(input, "first") {
		(*es.parent).changeStateByName("second")
		return "Exit with 1"
	} else
	if strings.Contains(input, "second") {
		(*es.parent).changeStateByName("second")
		return "Exit with 2"
	}

	return "Nothing done"
}

func (es *EntryState) GetParentScenario() *Scenario {
	return es.parent
}


type SecondState struct {
	parent *Scenario
}

func (ss *SecondState) InitWithParent(parent *Scenario) error {
	(*parent).registerState("second", ss)
}

func (ss *SecondState) RenderMessage() string {
	return "This is second message, you can only [exit] in order to get out of here"
}

func (ss *SecondState) HandleMessage(input string) string {
	if strings.Contains("exit", input) {
		(*ss.parent).changeStateByName("entry")
		return "Exiting..."
	}
	return "Not exit, stay here."
}

func (ss *SecondState) GetParentScenario() *Scenario {
	return ss.parent
}



func (rs *RootScenario) Name() string {
	return "RootScenario"
}
