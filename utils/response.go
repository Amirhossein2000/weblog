package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func BadResponseErr(w http.ResponseWriter) {
	responseBody := map[string]string{
		"message": "Bad request error",
	}
	WriteResponse(w, http.StatusBadRequest, responseBody)
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
		responseBody := map[string]string{
			"message": "Internal Server Error",
		}
		WriteResponse(w, http.StatusInternalServerError, responseBody)
	}
	w.WriteHeader(status)
	w.Write(byteBody)
}

func ArticleNotFound(w http.ResponseWriter) {
	WriteResponse(w, http.StatusNotFound, map[string]interface{}{
		"message": "Article does not exist",
	})
}
