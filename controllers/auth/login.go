package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"weblog/database"
	"weblog/schema"
	"weblog/utils"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	if utils.Authenticate(r) != nil {
		responseBody := map[string]string{
			"message": "already login",
		}
		utils.WriteResponse(w, http.StatusOK, responseBody)
	}

	loginRequest := schema.LoginRequest{}
	byteRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.InternalServerErr(w)
		return
	}

	err = json.Unmarshal(byteRequest, &loginRequest)
	if err != nil {
		log.Println("Unmarshal LoginRequest err:", err.Error())
		utils.BadJsonRequestStructure(w)
		return
	}

	user := schema.User{}
	err = database.DB.Where("email = ? AND password = ?",
		loginRequest.Email, loginRequest.Password).First(&user).Error
	if err != nil {
		log.Println("DB err:", err.Error())
		resp := map[string]string{
			"message": "wrong Email of Password",
		}
		utils.WriteResponse(w, http.StatusUnauthorized, resp)
		return
	}

	if &user == nil {
		resp := map[string]string{
			"message": "wrong Email of Password",
		}
		utils.WriteResponse(w, http.StatusUnauthorized, resp)

		return
	}

	resp := map[string]string{
		"message": "login successful",
		"token": utils.GenerateToken(
			&utils.AuthenticatedUser{
				Id:   user.ID,
				Role: user.Role,
			}),
	}
	utils.WriteResponse(w, http.StatusOK, resp)
}
