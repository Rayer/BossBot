package ChatBot

type ScenarioState interface {
	RenderMessage() string
	HandleMessage(input string) string
	GetParentScenario() *Scenario
}

