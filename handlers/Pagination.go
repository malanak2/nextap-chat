package handlers

import (
	"net/http"
	"strconv"
)

func GetVars(r *http.Request) (int, int, error) {
	var limit int
	var err error
	if r.URL.Query()["limit"] != nil {
		limit, err = strconv.Atoi(r.URL.Query()["limit"][0])
		if err != nil || limit < 1 || limit > 1000 {
			return 0, 0, ErrorInvalidLimit
		}
	}
	var page int
	if r.URL.Query()["page"] != nil {
		page, err = strconv.Atoi(r.URL.Query()["page"][0])
		if err != nil || page < 1 || page > 1000 {
			return 0, 0, ErrorInvalidPage
		}
	}
	return limit, page, nil
}
