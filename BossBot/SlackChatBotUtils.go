package BossBot

import (
	"ChatBot"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
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

func (ssi *SlackScenarioImpl) RenderSlackMessage() (string, []slack.Attachment, error) {
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

	msgHandler := s.KeywordHandler()
	attachment := msgHandler.GenerateAttachment(res)

	return res, []slack.Attachment{attachment}, nil
}

func (*SlackScenarioImpl) ResponseSlackCallbacks() {
	panic("implement me")
}

type SlackScenarioState interface {
	KeywordHandler() *KeywordHandler
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

type KeywordAction func(keyword string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error)

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
	kh.keywordList = append(kh.keywordList, *keyword)
}

func (kh *KeywordHandler) GenerateAttachment(input string) slack.Attachment {
	//Find all surrended by "[]"
	r, _ := regexp.Compile(`\[([A-Za-z 0-9_]*)]`)
	keywords := r.FindAllString(input, -1)

	var ret slack.Attachment
	var actions []slack.AttachmentAction

	ret.CallbackID = "chatbot-callback"

	for _, keywordDefine := range kh.keywordList {
		//TODO: Maybe we should use map to avoid O(n^2)?
		for _, keyword := range keywords {
			keyword = strings.Replace(keyword, "[", "", -1)
			keyword = strings.Replace(keyword, "]", "", -1)

			//TODO: Do we need case sensitive?
			if strings.ToLower(keywordDefine.Keyword) == strings.ToLower(keyword) {
				actions = append(actions, slack.AttachmentAction{
					Text:  strings.Title(keyword),
					Name:  strings.Title(keyword),
					Type:  "button",
					Value: keyword,
				})
				break
			}
		}
	}
	ret.Actions = actions
	return ret
}

func (kh *KeywordHandler) ParseAction(input string) (string, error) {
	for _, kw := range kh.keywordList {
		if strings.Contains(strings.ToLower(input), strings.ToLower(kw.Keyword)) {
			ret, err := kw.Action(kw.Keyword, kh.scenario, kh.state)
			if err != nil {
				return "", errors.Wrap(err, "Error parsing action : "+kw.Keyword)
			}
			return ret, nil
		}
	}
	return "No match keyword", nil
}
