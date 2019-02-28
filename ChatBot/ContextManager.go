package ChatBot

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type ContextManager struct {
	contextList   map[string]*UserContext
	Configuration *Configuration
}

type Configuration struct {
	ResetTimerSec int
}

func NewContextManager() *ContextManager {
	//Given a default value for configuration

	conf := Configuration{
		ResetTimerSec: 300,
	}

	return NewContextManagerWithConfig(&conf)
}

func NewContextManagerWithConfig(conf *Configuration) *ContextManager {
	ret := ContextManager{
		Configuration: conf,
	}
	ret.contextList = make(map[string]*UserContext)
	return &ret
}

func (cm *ContextManager) CreateUserContext(user string, entryScenario func() Scenario) *UserContext {
	uc := cm.contextList[user]
	if uc == nil {
		uc = NewUserContext(user, entryScenario())
		cm.contextList[user] = uc
	} else {
		log.Warnf("User context for %s already here, shouldn't try to get it first?", user)
	}
	return uc
}

func (cm *ContextManager) GetUserContext(user string) *UserContext {
	uc := cm.contextList[user]
	log.Debugf("Acception user : %s, current user list : %+v", user, uc)
	//Purge slice... it's stupid but it seems most maintainable way
	if uc != nil {
		log.Debugf("User %s, last session %v seconds ago...", user, time.Now().Sub(uc.LastAccess).Seconds())
	}
	if uc != nil && int(time.Now().Sub(uc.LastAccess).Seconds()) > cm.Configuration.ResetTimerSec {
		log.Infof("Re-Create ChatBot session %s due to timeout", user)
		cm.contextList[user] = nil
		return nil
	}
	return uc
}
