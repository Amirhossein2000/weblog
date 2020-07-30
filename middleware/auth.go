package middleware

import (
	"net/http"
	"weblog/utils"
)

type AuthRequireController func(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser)

func AuthRequire(w http.ResponseWriter, r *http.Request, handler AuthRequireController) {
	authUser := &utils.AuthenticatedUser{}
	authUser = utils.Authenticate(r)
	if authUser == nil {
		utils.PermissionDenied(w)
		return
	}

	handler(w, r, authUser)
}
