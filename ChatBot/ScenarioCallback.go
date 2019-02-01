package ChatBot

type ScenarioCallback interface {
	//Work like constructor
	InitScenario(uc *UserContext) error
	EnterScenario(source Scenario) error
	ExitScenario(askFrom Scenario) error
	DisposeScenario() error
}
