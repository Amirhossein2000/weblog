package comment

import (
	"log"
	"net/http"
	"weblog/database"
	"weblog/middleware"
	"weblog/schema"
	"weblog/utils"
)

func CommentController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
		return

	case http.MethodPost:
		middleware.AuthRequire(w, r, postHandler)
		return

	case http.MethodPut:
		middleware.AuthRequire(w, r, putHandler)
		return

	case http.MethodDelete:
		middleware.AuthRequire(w, r, deleteHandler)
		return

	default:
		utils.BadResponseErr(w)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.GetIDFromRequest(r)
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	comment := &schema.Comment{}
	err = database.DB.First(comment, articleID).Error
	if err != nil {
		responseBody := map[string]string{
			"message": "Comment not found",
		}
		utils.WriteResponse(w, http.StatusNotFound, responseBody)
		return
	}

	utils.WriteResponse(w, http.StatusOK, comment)
}

func postHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	comment, err := utils.UnmarshalComment(r.Body)
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	if !utils.ArticleExist(comment.ArticleID) {
		utils.ArticleNotFound(w)
	}

	comment.UserID = authUser.Id

	err = database.DB.Create(&comment).Error
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	responseBody := map[string]string{
		"message": "Created Successfully",
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}

func putHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	newComment, err := utils.UnmarshalComment(r.Body)
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	comment := schema.Comment{}
	database.DB.Select("user_id").Where("id = ?", newComment.ID).First(&comment)

	if !utils.HasPermission(comment.UserID, authUser) {
		utils.PermissionDenied(w)
		return
	}

	err = database.DB.Model(&comment).Where("id = ?", newComment.ID).
		UpdateColumn("value", newComment.Value).Error
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	responseBody := map[string]string{
		"message": "Updated Successfully",
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	commentID, err := utils.GetIDFromRequest(r)
	if err != nil {
		utils.BadResponseErr(w)
		return
	}

	comment := schema.Comment{}
	database.DB.Select("user_id").Where("id = ?", commentID).First(&comment)

	if !utils.HasPermission(comment.UserID, authUser) {
		utils.PermissionDenied(w)
		return
	}

	err = database.DB.Delete(&schema.Comment{}, commentID).Error
	if err != nil {
		log.Println("DB err:", err.Error())
	}

	responseBody := map[string]string{
		"message": "Deleted Successfully",
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}
