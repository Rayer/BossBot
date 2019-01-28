package ChatBot

type RootScenario struct{
	DefaultScenarioImpl
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {

}




func (rs *RootScenario) Name() string {
	return "RootScenario"
}
