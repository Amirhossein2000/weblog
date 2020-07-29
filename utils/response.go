package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func UnsupportedMethodErr(w http.ResponseWriter, method string) {
	BadResponseErrWithMessage(w, fmt.Sprintf("Method:%s is not supprted", method))
}

func BadJsonRequestStructure(w http.ResponseWriter) {
	BadResponseErrWithMessage(w, "Bad Json Request Structure")
}

func BadIDInURLErr(w http.ResponseWriter) {
	BadResponseErrWithMessage(w, "Bad ID in URL error, ID must be a number.")
}

func BadResponseErrWithMessage(w http.ResponseWriter, message string) {
	responseBody := map[string]string{
		"message": message,
	}
	WriteResponse(w, http.StatusBadRequest, responseBody)
}

func InternalServerErr(w http.ResponseWriter) {
	responseBody := map[string]string{
		"message": "Internal Server Error",
	}
	WriteResponse(w, http.StatusInternalServerError, responseBody)
}

func PermissionDenied(w http.ResponseWriter) {
	WriteResponse(w, http.StatusUnauthorized, map[string]string{
		"message": "Permission Denied",
	})
}

func WriteResponse(w http.ResponseWriter, status int, body interface{}) {
	byteBody, err := json.Marshal(body)
	if err != nil {
		log.Println("Response Unmarshal err:", err.Error())
		InternalServerErr(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(byteBody)
}

func ArticleNotFound(w http.ResponseWriter) {
	WriteResponse(w, http.StatusNotFound, map[string]interface{}{
		"message": "Article does not exist",
	})
}
