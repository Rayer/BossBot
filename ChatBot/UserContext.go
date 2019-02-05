package ChatBot

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

func NewUserContext(user string, rootScenario Scenario) *UserContext {
	ret := UserContext{
		user: user,
	}
	//Put root scenario into chain
	ret.scenarioChain = make([]Scenario, 0)
	ret.lastAccess = time.Now()
	rootScenario.SetUserContext(&ret)
	err := ret.InvokeNextScenario(rootScenario, Stack)
	if err != nil {
		log.Errorf("Error while trying to invoke root scenario : %s", err)
	}

	return &ret
}

func (uc *UserContext) GetCurrentScenario() Scenario {
	//TODO: should we check if there is NO root scenario?
	if len(uc.scenarioChain) == 0 {
		return nil
	}
	return uc.scenarioChain[len(uc.scenarioChain)-1]
}

func (uc *UserContext) GetRootScenario() Scenario {
	if len(uc.scenarioChain) == 0 {
		return nil
	}
	return uc.scenarioChain[0]
}

func (uc *UserContext) RenderMessage() (string, error) {
	uc.lastAccess = time.Now()
	ret, err := uc.GetCurrentScenario().RenderMessage()
	log.Infof("(%s)=>Rendering message : %s", uc.user, ret)
	return ret, err

}

func (uc *UserContext) HandleMessage(input string) (string, error) {
	uc.lastAccess = time.Now()
	ret, err := uc.GetCurrentScenario().HandleMessage(input)
	log.Infof("(%s)=>Rendering message : %s", uc.user, ret)
	return ret, err
}

func (uc *UserContext) InvokeNextScenario(scenario Scenario, strategy InvokeStrategy) error {

	thisScenario := uc.GetCurrentScenario()

	scenario.SetUserContext(uc)
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
		if oldScenario := uc.GetCurrentScenario(); oldScenario != nil {
			err := oldScenario.ExitScenario(scenario)
			if err != nil {
				log.Warnf("Error while exiting scenario '%s', error : %s", oldScenario.Name(), err)
			}
		}

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

func (uc *UserContext) ReturnLastScenario() error {
	var quitScenario Scenario
	var currentScenario Scenario
	quitScenario, uc.scenarioChain, currentScenario = uc.scenarioChain[len(uc.scenarioChain)-1], uc.scenarioChain[:len(uc.scenarioChain)-1], uc.scenarioChain[len(uc.scenarioChain)-1]

	err := quitScenario.ExitScenario(quitScenario)

	if err != nil {
		log.Warnf("Error while ExitScenario for %s, error : %s", quitScenario.Name(), err)
	}

	err = currentScenario.EnterScenario(quitScenario)

	if err != nil {
		log.Warnf("Error while EnterScenario for %s, error : %s", currentScenario.Name(), err)
	}

	err = quitScenario.DisposeScenario()

	if err != nil {
		log.Warnf("Error while DisposeScenario for %s, error : %s", quitScenario.Name(), err)
	}

	return nil
}
