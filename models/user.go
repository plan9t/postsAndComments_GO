package models

// User – пользователь в системе
type User struct {
	UserID    uint      `gorm:"column:user_id;primary_key"`
	FirstName string    `gorm:"size:32;not null"`
	LastName  string    `gorm:"size:32;not null"`
	Comments  []Comment // Связь один ко многим с комментариями
	Posts     []Post    // Связь один ко многим с постами
}
