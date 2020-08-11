package article

import (
	"fmt"
	"log"
	"net/http"
	"weblog/database"
	"weblog/middleware"
	"weblog/schema"
	"weblog/utils"
)

func ArticleController(w http.ResponseWriter, r *http.Request) {
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
		utils.UnsupportedMethodErr(w, r.Method)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.GetIDFromRequest(r)
	if err != nil {
		utils.BadIDInURLErr(w)
		return
	}

	article := &schema.Article{}
	err = database.DB.First(article, articleID).Error
	if err != nil {
		utils.ArticleNotFound(w)
		return
	}

	articleWithComments := &schema.ArticleWithComments{Article: article}
	parentComments := []schema.Comment{}

	err = database.DB.Find(&parentComments, "article_id = ? AND parent_comment_id = 0", article.ID).Error
	if err != nil {
		log.Println("DB err:", err.Error())
	}

	commentWithSubs := []*schema.CommentWithSubs{}
	if len(parentComments) > 0 {
		for i := 0; i < len(parentComments); i++ {
			commentWithSubs = append(commentWithSubs, utils.GetSubComments(parentComments[i]))
		}
		articleWithComments.Comments = commentWithSubs
	}

	utils.WriteResponse(w, http.StatusOK, articleWithComments)
}

func postHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	article, err := utils.UnmarshalArticle(r.Body)
	if err != nil {
		utils.BadJsonRequestStructure(w)
		return
	}

	article.UserID = authUser.Id

	err = database.DB.Create(&article).Error
	if err != nil {
		utils.InternalServerErr(w)
		return
	}

	responseBody := map[string]string{
		"message": fmt.Sprintf("Created Successfully, id = %d", article.ID),
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}

func putHandler(w http.ResponseWriter, r *http.Request, authUser *utils.AuthenticatedUser) {
	article, err := utils.UnmarshalArticle(r.Body)
	if err != nil {
		utils.BadJsonRequestStructure(w)
		return
	}

	oldArticle := schema.Article{}
	database.DB.Select("user_id").Where("id = ?", article.ID).First(&oldArticle)

	if !utils.HasPermission(oldArticle.UserID, authUser) {
		utils.PermissionDenied(w)
		return
	}

	article.UserID = oldArticle.UserID

	err = database.DB.Save(&article).Error
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
	articleID, err := utils.GetIDFromRequest(r)
	if err != nil {
		utils.BadIDInURLErr(w)
		return
	}

	article := schema.Article{}
	err = database.DB.Select("user_id").Where("id = ?", articleID).First(&article).Error
	if err != nil {
		log.Println("DB err:", err.Error())
	}

	if !utils.HasPermission(article.UserID, authUser) {
		utils.PermissionDenied(w)
		return
	}

	database.DB.Delete(&schema.Article{}, articleID)
	comments := []schema.Comment{}
	database.DB.Select("id").Find(&comments, "article_id = ?", articleID)

	utils.DeleteMultiCommentsWithReplys(comments)

	responseBody := map[string]string{
		"message": "Deleted Successfully",
	}
	utils.WriteResponse(w, http.StatusOK, responseBody)
}
