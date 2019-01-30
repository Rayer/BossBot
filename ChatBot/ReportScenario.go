package ChatBot

type ReportScenario struct {
	DefaultScenarioImpl
}

func (rs *ReportScenario) InitScenario(uc *UserContext) error {
	panic("implement me")
}

func (rs *ReportScenario) EnterScenario(source Scenario) error {
	panic("implement me")
}

func (rs *ReportScenario) ExitScenario(askFrom Scenario) error {
	panic("implement me")
}
