package models

import "gorm.io/gorm"

// User – пользователь в системе
type User struct {
	gorm.Model
	FirstName string    `gorm:"size:32;not null"`
	LastName  string    `gorm:"size:32;not null"`
	Comments  []Comment // Связь один ко многим с комментариями
	Posts     []Post    // Связь один ко многим с постами
}
