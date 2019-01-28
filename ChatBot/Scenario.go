package ChatBot

type Scenario interface {
	ScenarioCallback
	RenderMessage() string
	ProcessMessage(msg string) error
	GetUserContext() *UserContext
	Name() string
}

type DefaultScenarioImpl struct {

}

func (dsi *DefaultScenarioImpl) EnterScenario(source *Scenario) error {
	//Default action is : Raise state and log
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) ExitScenario(askFrom *Scenario) error {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) DisposeScenario() error {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) RenderMessage() string {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) ProcessMessage(msg string) error {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) GetUserContext() *UserContext {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) Name() string {
	panic("implement me")
}
