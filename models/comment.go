package models

import (
	"time"
)

// Comment - комментарий в посте
type Comment struct {
	CommentID       uint       `gorm:"column:comment_id;primary_key"`
	Content         string     `gorm:"size:2000;not null"`
	CreatedTime     time.Time  `gorm:"not null"`
	UserID          uint       // Внешний ключ для связи с пользователем
	User            User       // Связь многие к одному с пользователем
	PostID          uint       // Внешний ключ для связи с постом
	Post            Post       // Связь многие к одному с постом
	ParentCommentID *uint      // Внешний ключ для связи с родительским комментарием (опционально)
	ParentComment   *Comment   `gorm:"foreignKey:ParentCommentID"` // Связь многие к одному с родительским комментарием (опционально)
	ChildComments   []*Comment `gorm:"foreignKey:ParentCommentID"` // Связь один ко многим с дочерними комментариями
}
