package api

import (
	"encoding/json"
	"io/ioutil"
	"magnusquiz/pkg/log"
	"net/http"
	"strconv"
)

func createEntity(
	w http.ResponseWriter,
	r *http.Request,
	req interface{},
	saver func(interface{}) (string, error),
) {
	resp := new(HTTPResponse)
	resp.w = w
	if !resp.checkQkey(r) {
		resp.unauthorized()
		return
	}
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		resp.ise(err)
		return
	}
	_ = r.Body.Close()

	if err = json.Unmarshal(data, req); err != nil {
		resp.badRequest()
		return
	}

	uuid, err := saver(req)
	if err != nil {
		resp.ise(err)
		return
	}

	log.Logger.Infof("successfully create new entity with uuid : %s", uuid)

	resp.ok(map[string]string{"uuid": uuid})
}

func getIntFromQuery(r *http.Request, p string, dflt int) int {
	l, err := strconv.Atoi(r.URL.Query().Get(p))

	if err != nil || l <= 0 {
		log.Logger.Warnf(
			"err in converting param [%s], take default %d.ERR => %s\n", p, dflt, err,
		)
		l = dflt
	}
	return l
}
