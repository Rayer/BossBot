package ChatBot

type ScenarioState interface {
	RenderMessage() string
	HandleMessage() string
	GetParentScenario() *Scenario
}

