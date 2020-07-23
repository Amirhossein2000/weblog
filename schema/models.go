package schema

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"not null"`
	Username string `gorm:"not null"`
	Role     byte   `gorm:"not null"`
}

type Article struct {
	gorm.Model
	Title  string
	Body   string `gorm:"size:2000"`
	UserID uint
}

type Comment struct {
	gorm.Model
	ParentCommentID uint
	UserID          uint
	ArticleID       uint
	Value           string
}
