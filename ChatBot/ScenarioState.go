package ChatBot

type ScenarioState interface {
	InitWithParent(parent Scenario) error
	RenderMessage() (string, error)
	HandleMessage(input string) (string, error)
	GetParentScenario() Scenario
}
