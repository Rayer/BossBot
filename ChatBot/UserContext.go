package ChatBot

import (
	"github.com/pkg/errors"
	"time"
)

type UserContext struct {
	user          string
	scenarioChain []Scenario
	lastAccess    time.Time
}

type InvokeStrategy int

const (
	Stack   InvokeStrategy = 0
	Trim    InvokeStrategy = 1
	Replace InvokeStrategy = 2
)

func NewUserContext(user string) *UserContext {
	ret := UserContext{
		user: user,
	}
	//Put root scenario into chain
	rs := RootScenario{}
	rs.stateList = make(map[string]ScenarioState)
	ret.InvokeNextScenario(&rs, Stack)

	return &ret
}

func (uc *UserContext) GetCurrentScenario() Scenario {
	//TODO: should we check if there is NO root scenario?
	if len(uc.scenarioChain) == 0 {
		return nil
	}
	return uc.scenarioChain[len(uc.scenarioChain)-1]
}

func (uc *UserContext) RenderMessage() (string, error) {
	return uc.GetCurrentScenario().RenderMessage()
}

func (uc *UserContext) HandleMessage(input string) (string, error) {
	return uc.GetCurrentScenario().HandleMessage(input)
}

func (uc *UserContext) InvokeNextScenario(scenario Scenario, strategy InvokeStrategy) error {

	thisScenario := uc.GetCurrentScenario()

	err := scenario.InitScenario(uc)

	if err != nil {
		return errors.Wrap(err, "Fail to init scenario : "+scenario.Name())
	}

	err = scenario.EnterScenario(thisScenario)

	if err != nil {
		return errors.Wrap(err, "Fail to enter scenario : "+scenario.Name())
	}
	switch strategy {
	case Stack:
		uc.scenarioChain = append(uc.scenarioChain, scenario)
	case Trim:
		//Remove from 1 to end of slice
		for idx, s := range uc.scenarioChain {
			if idx == 0 {
				continue
			}
			err = s.ExitScenario(thisScenario)
			if err != nil {
				return errors.Wrap(err, "Error while exiting scenario : "+s.Name())
			}
		}
		uc.scenarioChain = append([]Scenario{}, uc.scenarioChain[0], scenario)

	case Replace:
		//TODO: Root scenario can't be replaced
		old := uc.scenarioChain[len(uc.scenarioChain)-1]
		err = old.ExitScenario(thisScenario)
		if err != nil {
			return errors.Wrap(err, "Error while exiting scenario : "+old.Name())
		}
		uc.scenarioChain[len(uc.scenarioChain)-1] = thisScenario
	}
	return nil
}
