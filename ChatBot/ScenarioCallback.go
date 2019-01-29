package ChatBot

type ScenarioCallback interface {
	//Work like constructor
	EnterScenario(source Scenario) error
	ExitScenario(askFrom Scenario) error
	DisposeScenario() error
}
