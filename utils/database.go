package utils

import (
	"weblog/database"
	"weblog/schema"
)

func GetSubComments(comment schema.Comment) *schema.CommentWithSubs {
	output := &schema.CommentWithSubs{Comment: &comment}
	subComments := []schema.Comment{}
	database.DB.Where("parent_comment_id = ?", comment.ID).Find(&subComments)
	if len(subComments) > 0 {
		for i := 0; i < len(subComments); i++ {
			output.Replys = append(output.Replys, GetSubComments(subComments[i]))
		}
	}

	return output
}

func ArticleExist(articleID uint) bool {
	article := schema.Article{}
	database.DB.Select("id").Where("id = ?", articleID).First(&article)
	if article.ID != 0 {
		return true
	}

	return false
}
