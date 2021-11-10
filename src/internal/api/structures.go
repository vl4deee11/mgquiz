package api

import (
	"magnusquiz/pkg/enum"
)

type createQuestionReq struct {
	Text            string               `json:"text"`
	DifficultyLevel enum.DifficultyLevel `json:"difficulty_level"`
}

type createAnswerReq struct {
	IsRight      bool   `json:"is_right"`
	Text         string `json:"text"`
	QuestionUUID string `json:"question_uuid"`
}

type createUserInfoReq struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Link     string `json:"link"`
	Question string `json:"question"`
}
