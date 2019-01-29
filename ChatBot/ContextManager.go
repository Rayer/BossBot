package ChatBot

type ContextManager struct {
	contextList map[string]*UserContext
}

func NewContextManager() *ContextManager{
	ret := ContextManager{}
	ret.contextList = make(map[string]*UserContext)
	return &ret
}

func (cm *ContextManager) GetUserContext(user string) *UserContext {
	uc := cm.contextList[user]
	if uc == nil {
		rs := new RootScenario{}
		uc = &UserContext{
			user: user,
			scenarioChain:[]*Scenario{rs.(*Scenario)},
		}
	}



}