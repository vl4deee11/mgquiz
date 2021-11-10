package api

import (
	"magnusquiz/pkg/repo"
	"net/http"
)

func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	saver := func(raw interface{}) (string, error) {
		req := raw.(*createAnswerReq)
		uuid, err := repo.NewAnswer(req.Text, req.QuestionUUID, req.IsRight)
		if err != nil {
			return "", err
		}

		return uuid, nil
	}

	createEntity(w, r, &createAnswerReq{}, saver)
}
