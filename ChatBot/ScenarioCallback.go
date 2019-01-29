package ChatBot

type ScenarioCallback interface {
	EnterScenario(source Scenario) error
	ExitScenario(askFrom Scenario) error
	DisposeScenario() error
}
