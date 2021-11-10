package repo

type RAnswer struct {
	IsRight bool   `json:"is_right"`
	Text    string `json:"text"`
	UUID    string `json:"uuid"`
}

type RQuestion struct {
	HasRightAnswer bool      `json:"-"`
	Text           string    `json:"text"`
	Answers        []RAnswer `json:"answers"`
}
