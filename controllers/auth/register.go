package auth

import (
	"net/http"
	"weblog/database"
	"weblog/schema"
	"weblog/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if utils.Authenticate(r) != nil {
		responseBody := map[string]string{
			"message": "already login",
		}
		utils.WriteResponse(w, http.StatusOK, responseBody)
	}

	user, err := utils.UnmarshalUser(r.Body)
	if err != nil {
		utils.BadJsonRequestStructure(w)
		return
	}

	user.Role = utils.UserRole

	oldUser := &schema.User{}
	database.DB.Select("email").Where("email = ?", user.Email).First(oldUser)
	if oldUser.Email != "" {
		responseBody := map[string]string{
			"message": "email has been taken",
		}
		utils.WriteResponse(w, http.StatusOK, responseBody)
		return
	}

	database.DB.Select("username").Where("username = ?", user.Username).First(oldUser)
	if oldUser.Username != "" {
		responseBody := map[string]string{
			"message": "username has been taken",
		}
		utils.WriteResponse(w, http.StatusOK, responseBody)
		return
	}

	err = database.DB.Create(&user).Error
	if err != nil || user == nil {
		utils.InternalServerErr(w)
		return
	}

	responseBody := map[string]string{
		"message": "Register Successfully",
		"token": utils.GenerateToken(
			&utils.AuthenticatedUser{
				Id:   user.ID,
				Role: user.Role,
			}),
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}
