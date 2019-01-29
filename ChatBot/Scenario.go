package ChatBot

import log "github.com/sirupsen/logrus"

type Scenario interface {
	ScenarioCallback
	RenderMessage() string
	ProcessMessage(msg string) error
	GetUserContext() *UserContext
	Name() string

	getState(name string) *ScenarioState
	changeStateByName(name string) error
	registerState(name string, state interface{})
}

type DefaultScenarioImpl struct {
	stateList map[string]*ScenarioState
	currentState *ScenarioState
}

func (dsi *DefaultScenarioImpl) getState(name string) *ScenarioState {
	return dsi.stateList[name]
}

func (dsi *DefaultScenarioImpl) changeStateByName(name string) error {
	state := dsi.stateList[name]
	dsi.currentState = state
	return nil
}

func (dsi *DefaultScenarioImpl) registerState(name string, state interface{}) {
	stateImpl := state.(*ScenarioState)
	dsi.stateList[name] = stateImpl
}

func (dsi *DefaultScenarioImpl) EnterScenario(source *Scenario) error {
	log.Debugln("Entering scenario")
	return nil
}

func (dsi *DefaultScenarioImpl) ExitScenario(askFrom *Scenario) error {
	log.Debugln("Exiting scenario")
	return nil
}

func (dsi *DefaultScenarioImpl) DisposeScenario() error {
	log.Debugln("Disposing Scenario")
	return nil
}

func (dsi *DefaultScenarioImpl) RenderMessage() string {
	panic("implement me")
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

