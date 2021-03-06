package assistant

import (
	"github.com/COSCUP/assistant/program-fetcher"
	log "github.com/Sirupsen/logrus"
)

type AskProgramByProgramIntentProcessor struct {
}

func (AskProgramByProgramIntentProcessor) Name() string {
	// return "projects/coscup/agent/intents/30d89d56-6dac-4649-a940-71c99eb69324"
	return "Intent Ask Program by Program"
}

func (p AskProgramByProgramIntentProcessor) displayMessage(sessionTitle string) string {
	return "「" + sessionTitle + "」的議程資訊如下："
}

func (p AskProgramByProgramIntentProcessor) speechMessage(sessionTitle string) string {
	return "議程資訊如下"
}

func (p AskProgramByProgramIntentProcessor) getSuggsetion(inFavoriteList bool, session *fetcher.Session) []map[string]interface{} {
	return getSuggestionWithSession(inFavoriteList, session)
}
func getSuggestionWithSession(inFavoriteList bool, session *fetcher.Session) []map[string]interface{} {
	ret := []map[string]interface{}{}

	if !inFavoriteList {
		ret = append(ret, getSuggestionPayload("🌟我有興趣"))
	} else {
		ret = append(ret, getSuggestionPayload("移除議程"))
	}

	ret = append(ret, getSuggestionPayload(session.Room+"在哪"))
	ret = append(ret, getSuggestionPayload(session.Room+"的下一場議程什麼時候開始"))
	dt := "第一天"
	if IsDayTwo(session.Start) {
		dt = "第二天"
	}
	timeLine := dt + session.End.Format("15:04")
	ret = append(ret, getSuggestionPayload(timeLine+"之後有哪些議程"))
	return ret
}
func (p AskProgramByProgramIntentProcessor) Payload(input *DialogflowRequest) map[string]interface{} {
	perviousDisplayedSessionListInfo := input.Context("pervious_session_list")
	log.Println("perviousDisplayedSessionList:", perviousDisplayedSessionListInfo)
	number := input.SelectedNumber()
	selectedID := ""
	list := perviousDisplayedSessionListInfo["list"].([]interface{})

	if number >= 1 && len(list) > number-1 {
		//
		selectedID = list[number-1].(string)
	}

	prog, _ := fetcher.GetPrograms()
	sessionInfo := prog.GetSessionByID(selectedID)
	title := sessionInfo.Zh.Title
	desc := sessionInfo.Zh.Description
	timeLine := sessionInfo.Start.Format("15:04") + "~" + sessionInfo.End.Format("15:04")
	subTitle := sessionInfo.Room + " " + timeLine

	sessionPhotoUrl := sessionInfo.SpeakerPhotoUrl()

	inFavoriteList := NewUserStorageFromDialogflowRequest(input).isSessionIdInFavorite(selectedID)

	return map[string]interface{}{
		"expectUserResponse": true,

		// "systemIntent": getListSystemIntentPayload(),
		"richResponse": map[string]interface{}{
			"items": []map[string]interface{}{
				getSimpleResponsePayload(p.speechMessage(title), p.displayMessage(title)),
				getBasicCardResponsePayload(
					title,
					subTitle,
					desc,
					sessionPhotoUrl, "講者照片",
					"議程網頁", "https://coscup.org/2020/zh-TW/agenda/"+sessionInfo.ID, "CROPPED"),

				// getSimpleResponsePayload("123", "321"),
				// getTableCardResponsePayload("title", "subtitle",
				// 	[]Row{
				// 		getRowPayload([]Cell{getCellPayload("1"), getCellPayload("2"), getCellPayload("3")}, true),
				// 		getRowPayload([]Cell{getCellPayload("12"), getCellPayload("23"), getCellPayload("34")}, true),
				// 	},
				// 	[]ColunmProperty{
				// 		getColumnPropertyPayload("C1", HorizontalAlignmentCenter),
				// 		getColumnPropertyPayload("C2", HorizontalAlignmentCenter),
				// 		getColumnPropertyPayload("C3", HorizontalAlignmentCenter),
				// 	},
				// 	"https://coscup.org/2019/_nuxt/img/c2f9236.png", "image", "按鈕", "https://www.tih.tw", "CROPPED",
				// ),
			},
			"suggestions": p.getSuggsetion(inFavoriteList, sessionInfo),

			// "linkOutSuggestion": getLinkOutSuggestionPayload("tih", "https://www.tih.tw"),
		},
		"outputContexts": map[string]interface{}{
			"selected_session": map[string]interface{}{
				"id": selectedID,
			},
		},
	}
}
