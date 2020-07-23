package timeline

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"weblog/database"
	"weblog/schema"
	"weblog/utils"
)

func TimelineController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		postHandler(w, r)
		return

	default:
		utils.BadResponseErr(w)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	timelineRequest := schema.TimelineRequest{}
	byteRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	err = json.Unmarshal(byteRequest, &timelineRequest)
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	articleIds := []uint{}
	articles := []schema.Article{}

	if timelineRequest.UserId == 0 {
		err = database.DB.Select("id").Offset(timelineRequest.Offset).
			Limit(timelineRequest.Limit).Find(&articles).Error
	} else {
		err = database.DB.Select("id").Offset(timelineRequest.Offset).
			Limit(timelineRequest.Limit).Where("user_id = ?", timelineRequest.UserId).Find(&articles).Error
	}

	if err != nil {
		log.Println("DB err:", err.Error())
		utils.BadResponseErr(w)
		return
	}

	if len(articles) > 0 {
		for i := 0; i < len(articles); i++ {
			articleIds = append(articleIds, articles[i].ID)
		}
	}

	responseBody := map[string]interface{}{
		"message": articleIds,
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}
