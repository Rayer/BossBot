package BossBot

import (
	"ChatBot"
	"github.com/nlopes/slack"
	"regexp"
	"strings"
)

type SlackScenario interface {
	ChatBot.Scenario
	RenderSlackMessage() (string, []slack.Attachment, error)
	ResponseSlackCallbacks()
}

type SlackScenarioImpl struct {
	parentScenario ChatBot.Scenario
}

func (ssi *SlackScenarioImpl) InitSlackScenario(scenario ChatBot.Scenario) {
	ssi.parentScenario = scenario
}

func (ssi *SlackScenarioImpl) RenderSlackMessage(input string) (string, []slack.Attachment, error) {
	currentState := ssi.parentScenario.GetCurrentState()
	res, err := currentState.RenderMessage()
	if err != nil {
		//log.Errorf(err)
		return "Error!", nil, err
	}

	s, ok := currentState.(SlackScenarioState)

	if !ok {
		return res, nil, nil
	}

	msgHandler := s.GetKeywordHandler()
	attachment := msgHandler.GenerateAttachment(res)

	return res, []slack.Attachment{attachment}, nil
}

func (*SlackScenarioImpl) ResponseSlackCallbacks() {
	panic("implement me")
}

type SlackScenarioState interface {
	GetKeywordHandler() *KeywordHandler
}

type SlackScenarioStateImpl struct {
	keywordHandler *KeywordHandler
}

func (s *SlackScenarioStateImpl) KeywordHandler() *KeywordHandler {
	return s.keywordHandler
}

func NewSlackScenarioStateImpl(state ChatBot.ScenarioState) *SlackScenarioStateImpl {
	return &SlackScenarioStateImpl{keywordHandler: NewKeywordHandler(state.GetParentScenario(), state)}
}

type KeywordAction func(keyword string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) string

type Keyword struct {
	Keyword string
	Action  KeywordAction
}

type KeywordHandler struct {
	keywordList []Keyword
	scenario    ChatBot.Scenario
	state       ChatBot.ScenarioState
}

func NewKeywordHandler(scenario ChatBot.Scenario, state ChatBot.ScenarioState) *KeywordHandler {
	return &KeywordHandler{scenario: scenario, state: state}
}

func (kh *KeywordHandler) RegisterKeyword(keyword *Keyword) {
	if kh.keywordList == nil {
		kh.keywordList = []Keyword{}
	}
}

func (kh *KeywordHandler) GenerateAttachment(input string) slack.Attachment {
	//Find all surrended by "[]"
	r, _ := regexp.Compile(`\[([A-Za-z 0-9_]*)]`)
	keywords := r.FindAllString(input, -1)

	var ret slack.Attachment
	var actions []slack.AttachmentAction

	for _, keywordDefine := range kh.keywordList {
		//TODO: Maybe we should use map to avoid O(n^2)?
		for _, keyword := range keywords {
			if keywordDefine.Keyword == keyword {
				actions = append(actions, slack.AttachmentAction{
					Text: strings.Title(keyword),
					Type: "button",
				})
				break
			}
		}
	}
	ret.Actions = actions
	return ret
}

func (kh *KeywordHandler) ParseAction(input string) error {
	return nil
}
