package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetIDFromRequest(r *http.Request) (uint, error) {
	strID := r.URL.Query().Get("id")
	id, err := strconv.Atoi(strID)
	if err != nil || strID == "" || id < 1 {
		return 0, fmt.Errorf("Bad id: %s", strID)
	}
	return uint(id), nil
}
