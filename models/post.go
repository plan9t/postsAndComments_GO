package models

import "time"

// Post представляет пост пользователя в системе.
type Post struct {
	ID          int       // Уникальный идентификатор
	Title       string    // Заголовок поста
	Content     string    // Содержимое поста
	Commentable bool      // Флаг, указывающий, можно ли комментировать пост
	CreatedTime time.Time // Время создания поста
	UserID      int       // ID пользователя, создавшего пост
}
