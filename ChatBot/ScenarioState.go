package ChatBot

type ScenarioState interface {
	InitScenarioState(scenario Scenario)
	RenderMessage() (string, error)
	HandleMessage(input string) (string, error)
	GetParentScenario() Scenario
	SetParentScenario(parent Scenario)
}

type DefaultScenarioStateImpl struct {
	parent Scenario
}

func (dssi *DefaultScenarioStateImpl) GetParentScenario() Scenario {
	return dssi.parent
}

func (dssi *DefaultScenarioStateImpl) SetParentScenario(parent Scenario) {
	dssi.parent = parent
}
