package ChatBot

type Scenario interface {
	ScenarioCallback
	RenderMessage() (string, error)
	HandleMessage(input string) (string, error)
	SetUserContext(user *UserContext)
	GetUserContext() *UserContext
	Name() string

	getState(name string) ScenarioState
	changeStateByName(name string) error
	registerState(name string, state ScenarioState, parentScenario Scenario)
}

type DefaultScenarioImpl struct {
	stateList    map[string]ScenarioState
	currentState ScenarioState
	userContext  *UserContext
}

func (dsi *DefaultScenarioImpl) getState(name string) ScenarioState {
	return dsi.stateList[name]
}

func (dsi *DefaultScenarioImpl) changeStateByName(name string) error {
	state := dsi.stateList[name]
	if state == nil {
		panic("Can't find state " + name + " in the scenario!")
	}
	dsi.currentState = state
	return nil
}

func (dsi *DefaultScenarioImpl) registerState(name string, state ScenarioState, parentScenario Scenario) {
	state.SetParentScenario(parentScenario)
	dsi.stateList[name] = state
	if dsi.currentState == nil {
		dsi.currentState = state
	}
}

func (dsi *DefaultScenarioImpl) RenderMessage() (string, error) {
	return dsi.currentState.RenderMessage()
}

func (dsi *DefaultScenarioImpl) HandleMessage(input string) (string, error) {
	return dsi.currentState.HandleMessage(input)
}

func (dsi *DefaultScenarioImpl) SetUserContext(user *UserContext) {
	dsi.userContext = user
}

func (dsi *DefaultScenarioImpl) GetUserContext() *UserContext {
	return dsi.userContext
}
