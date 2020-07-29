package user

import (
	"net/http"
	"weblog/database"
	"weblog/middleware"
	"weblog/schema"
	"weblog/utils"
)

func UserController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
		return

	case http.MethodPut:
		middleware.AuthRequire(w, r, putHandler)
		return

	case http.MethodDelete:
		middleware.AuthRequire(w, r, deleteHandler)
		return

	default:
		utils.UnsupportedMethodErr(w, r.Method)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetIDFromRequest(r)
	if err != nil {
		utils.BadIDInURLErr(w)
		return
	}

	user := &schema.User{}
	err = database.DB.First(user, userID).Error
	if err != nil {
		responseBody := map[string]string{
			"message": "User not found",
		}
		utils.WriteResponse(w, http.StatusNotFound, responseBody)
		return
	}
	utils.WriteResponse(w, http.StatusOK, user)
}

func putHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	updateUser, err := utils.UnmarshalUser(r.Body)
	if err != nil {
		utils.BadJsonRequestStructure(w)
		return
	}

	if !utils.HasPermission(updateUser.ID, authUser) {
		utils.PermissionDenied(w)
		return
	}

	user := schema.User{}
	database.DB.First(&user, updateUser.ID)

	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	if updateUser.Password != "" {
		user.Password = updateUser.Password
	}

	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}

	err = database.DB.Save(&updateUser).Error
	if err != nil {
		utils.InternalServerErr(w)
		return
	}

	responseBody := map[string]string{
		"message": "Updated Successfully",
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	userId, err := utils.GetIDFromRequest(r)
	if err != nil {
		utils.BadIDInURLErr(w)
		return
	}

	if !utils.HasPermission(userId, authUser) {
		utils.PermissionDenied(w)
		return
	}

	database.DB.Delete(&schema.Comment{}, "user_id = ?", userId)
	database.DB.Delete(&schema.Article{}, "user_id = ?", userId)
	database.DB.Delete(&schema.User{}, userId)

	responseBody := map[string]string{
		"message": "Deleted Successfully",
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}
