package ChatBot

type ScenarioState interface {
	InitWithParent(parent Scenario) error
	RenderMessage() string
	HandleMessage(input string) string
	GetParentScenario() Scenario
}
