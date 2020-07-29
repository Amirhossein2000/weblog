package auth

import (
	"net/http"
	"weblog/database"
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

	err = database.DB.Create(&user).Error
	if err != nil || user == nil {
		responseBody := map[string]string{
			"message": "username or email has been taken",
		}
		utils.WriteResponse(w, http.StatusOK, responseBody)
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
