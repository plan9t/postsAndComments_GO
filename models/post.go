package models

import (
	"time"
)

// Post представляет пост пользователя в системе.
type Post struct {
	PostID      uint      `gorm:"primary_key"`
	Title       string    `gorm:"size:64;not null"`
	Content     string    `gorm:"type:text;not null"`
	Commentable bool      `gorm:"not null"`
	CreatedTime time.Time `gorm:"not null"`
	UserID      uint      // Внешний ключ для связи с пользователем
	User        User      // Связь многие к одному с пользователем
	Comments    []Comment // Связь один ко многим с комментариями
}
