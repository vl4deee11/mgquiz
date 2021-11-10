package api

import (
	"magnusquiz/pkg/repo"
	"net/http"
)

func GenerateQuestions(w http.ResponseWriter, r *http.Request) {
	resp := new(HTTPResponse)
	resp.w = w
	n := getIntFromQuery(r, "n", 9)
	m, err := repo.GetNRandomQuestions(n)
	if err != nil {
		resp.ise(err)
		return
	}

	res := make([]*repo.RQuestion, 0, len(m))
	for k := range m {
		res = append(res, m[k])
	}

	resp.ok(res)
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	saver := func(raw interface{}) (string, error) {
		req := raw.(*createQuestionReq)
		uuid, err := repo.NewQuestion(req.Text, req.DifficultyLevel)
		if err != nil {
			return "", err
		}

		return uuid, nil
	}

	createEntity(w, r, &createQuestionReq{}, saver)
}
