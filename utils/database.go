package utils

import (
	"weblog/database"
	"weblog/schema"
)

func GetSubComments(comment schema.Comment) *schema.CommentWithSubs {
	output := &schema.CommentWithSubs{Comment: &comment}
	replys := []schema.Comment{}
	database.DB.Find(&replys, "parent_comment_id = ?", comment.ID)
	if len(replys) > 0 {
		for i := 0; i < len(replys); i++ {
			output.Replys = append(output.Replys, GetSubComments(replys[i]))
		}
	}

	return output
}

func DeleteSubComments(commentId uint) {
	database.DB.Delete(&schema.Comment{}, "id = ?", commentId)
	replys := []schema.Comment{}

	database.DB.Select("id").Find(&replys, "parent_comment_id = ?", commentId)

	DeleteMultiCommentsWithReplys(replys)
}

func DeleteMultiCommentsWithReplys(comments []schema.Comment) {
	if len(comments) == 0 {
		return
	}

	for i := 0; i < len(comments); i++ {
		DeleteSubComments(comments[i].ID)
	}
}

func ArticleExist(articleID uint) bool {
	article := schema.Article{}
	database.DB.Select("id").Where("id = ?", articleID).First(&article)
	if article.ID != 0 {
		return true
	}

	return false
}
