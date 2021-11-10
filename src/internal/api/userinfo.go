package api

import (
	"magnusquiz/pkg/repo"
	"net/http"
)

func CreateUserInfo(w http.ResponseWriter, r *http.Request) {
	saver := func(raw interface{}) (string, error) {
		req := raw.(*createUserInfoReq)
		uuid, err := repo.NewUserInfo(
			req.Name,
			req.Phone,
			req.Email,
			req.Link,
			req.Question,
		)
		if err != nil {
			return "", err
		}

		return uuid, nil
	}

	createEntity(w, r, &createUserInfoReq{}, saver)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	resp := new(HTTPResponse)
	resp.w = w
	if !resp.checkQkey(r) {
		resp.unauthorized()
		return
	}

	l := getIntFromQuery(r, "l", 100)
	o := getIntFromQuery(r, "o", 0)

	res, err := repo.GetUserInfoList(l, o)
	if err != nil {
		resp.ise(err)
		return
	}

	resp.ok(res)
}
