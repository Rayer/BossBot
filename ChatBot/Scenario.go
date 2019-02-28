package ChatBot

type Scenario interface {
	ScenarioCallback
	RenderMessage() (string, error)
	HandleMessage(input string) (string, error)
	SetUserContext(user *UserContext)
	GetUserContext() *UserContext
	Name() string

	GetState(name string) ScenarioState
	GetCurrentState() ScenarioState
	ChangeStateByName(name string) error
	RegisterState(name string, state ScenarioState, parentScenario Scenario)
}

type DefaultScenarioImpl struct {
	stateList    map[string]ScenarioState
	currentState ScenarioState
	userContext  *UserContext
}

func (dsi *DefaultScenarioImpl) GetCurrentState() ScenarioState {
	return dsi.currentState
}

func (dsi *DefaultScenarioImpl) InitScenario(uc *UserContext) {
	dsi.stateList = make(map[string]ScenarioState)
	dsi.userContext = uc
}

func (dsi *DefaultScenarioImpl) GetState(name string) ScenarioState {
	return dsi.stateList[name]
}

func (dsi *DefaultScenarioImpl) ChangeStateByName(name string) error {
	state := dsi.stateList[name]
	if state == nil {
		panic("Can't find state " + name + " in the scenario!")
	}
	dsi.currentState = state
	return nil
}

func (dsi *DefaultScenarioImpl) RegisterState(name string, state ScenarioState, parentScenario Scenario) {
	state.SetParentScenario(parentScenario)
	state.InitScenarioState(parentScenario)
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
