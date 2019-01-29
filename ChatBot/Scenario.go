package ChatBot

import log "github.com/sirupsen/logrus"

type Scenario interface {
	ScenarioCallback
	RenderMessage() string
	HandleMessage(input string) string
	ProcessMessage(msg string) error
	GetUserContext() *UserContext
	Name() string

	getState(name string) ScenarioState
	changeStateByName(name string) error
	registerState(name string, state ScenarioState)
}

type DefaultScenarioImpl struct {
	stateList    map[string]ScenarioState
	currentState ScenarioState
}

func (dsi *DefaultScenarioImpl) getState(name string) ScenarioState {
	return dsi.stateList[name]
}

func (dsi *DefaultScenarioImpl) changeStateByName(name string) error {
	state := dsi.stateList[name]
	dsi.currentState = state
	return nil
}

func (dsi *DefaultScenarioImpl) registerState(name string, state ScenarioState) {
	dsi.stateList[name] = state
	if dsi.currentState == nil {
		dsi.currentState = state
	}
}

func (dsi *DefaultScenarioImpl) EnterScenario(source Scenario) error {
	log.Debugln("Entering scenario")
	return nil
}

func (dsi *DefaultScenarioImpl) ExitScenario(askFrom Scenario) error {
	log.Debugln("Exiting scenario")
	return nil
}

func (dsi *DefaultScenarioImpl) DisposeScenario() error {
	log.Debugln("Disposing Scenario")
	return nil
}

func (dsi *DefaultScenarioImpl) RenderMessage() string {
	return dsi.currentState.RenderMessage()
}

func (dsi *DefaultScenarioImpl) HandleMessage(input string) string {
	return dsi.currentState.HandleMessage(input)
}

//With keyword system we can have a default process message. However, not now.
func (dsi *DefaultScenarioImpl) ProcessMessage(msg string) error {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) GetUserContext() *UserContext {
	panic("implement me")
}

func (dsi *DefaultScenarioImpl) Name() string {
	panic("implement me")
}
