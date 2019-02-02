package ChatBot

type ContextManager struct {
	contextList map[string]*UserContext
}

func NewContextManager() *ContextManager {
	ret := ContextManager{}
	ret.contextList = make(map[string]*UserContext)
	return &ret
}

func (cm *ContextManager) CreateUserContext(user string, entryScenario func() Scenario) *UserContext {
	uc := cm.contextList[user]
	if uc == nil {
		uc = NewUserContext(user, entryScenario())
		cm.contextList[user] = uc
	}
	return uc

}

func (cm *ContextManager) GetUserContext(user string) *UserContext {
	uc := cm.contextList[user]
	return uc
}
