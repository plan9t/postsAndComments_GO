package main

import (
	"OZON/models"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"log"
)

var (
	db           *gorm.DB
	userType     *graphql.Object
	postType     *graphql.Object
	commentType  *graphql.Object
	queryType    *graphql.Object
	mutationType *graphql.Object
)

func initGraphQLTypes(db *gorm.DB) {
	commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: (graphql.FieldsThunk)(func() graphql.Fields { // Через замыкание
				return graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.Int,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							log.Printf("Type of p.Source: %T", p.Source)   // Проверка типа p.Source
							log.Printf("Value of p.Source: %#v", p.Source) // Вывод значения p.Source
							if post, ok := p.Source.(*models.Comment); ok {
								log.Printf("Retrieved comment ID: %v", post.PostID)
								return post.CommentID, nil
							} else {
								log.Printf("p.Source is not of type *models.Comment")
							}
							return nil, fmt.Errorf("could not retrieve comment ID")
						},
					},
					"content": &graphql.Field{
						Type: graphql.String,
					},
					"createdTime": &graphql.Field{
						Type: graphql.DateTime,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if comment, ok := p.Source.(*models.Comment); ok && comment != nil {
								log.Printf("Retrieving created time for comment ID: %v", comment.CommentID)
								return comment.CreatedTime, nil
							}
							log.Printf("Failed to retrieve created time: comment not found or incorrect type")
							return nil, fmt.Errorf("could not retrieve created time for comment")
						},
					},
					"user": &graphql.Field{
						Type: graphql.Int,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if comment, ok := p.Source.(*models.Comment); ok {
								// Возвращаем ID пользователя, связанного с комментарием
								return comment.UserID, nil
							}
							log.Printf("Failed to retrieve user ID: comment not found or incorrect type")
							return nil, fmt.Errorf("could not retrieve user ID for comment")
						},
					},
					"parent_id": &graphql.Field{ // Добавляем новое поле parent_id
						Type: graphql.Int,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if comment, ok := p.Source.(*models.Comment); ok {
								return comment.ParentCommentID, nil // Убедитесь, что ParentID определен в модели Comment
							}
							return nil, fmt.Errorf("could not retrieve parent ID for comment")
						},
					},
					"children": &graphql.Field{
						Type: graphql.NewList(commentType),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if comment, ok := p.Source.(*models.Comment); ok {
								return getNestedComments(db, comment.CommentID)
							}
							return nil, fmt.Errorf("could not retrieve children comments")
						},
					},
				}
			}),
		},
	)

	postType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Post",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						log.Printf("Type of p.Source: %T", p.Source)   // Проверка типа p.Source
						log.Printf("Value of p.Source: %#v", p.Source) // Вывод значения p.Source
						if post, ok := p.Source.(*models.Post); ok {
							log.Printf("Retrieved post ID: %v", post.PostID)
							return post.PostID, nil
						} else {
							log.Printf("p.Source is not of type *models.Post")
						}
						return nil, fmt.Errorf("could not retrieve post ID")
					},
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"content": &graphql.Field{
					Type: graphql.String,
				},
				"commentable": &graphql.Field{
					Type: graphql.Boolean,
				},
				"createdTime": &graphql.Field{
					Type: graphql.DateTime,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
					Args: graphql.FieldConfigArgument{
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"offset": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if post, ok := p.Source.(*models.Post); ok {
							var comments []*models.Comment
							query := db.Where("post_id = ?", post.PostID)
							if limit, ok := p.Args["limit"].(int); ok {
								query = query.Limit(limit)
							}
							if offset, ok := p.Args["offset"].(int); ok {
								query = query.Offset(offset)
							}
							result := query.Find(&comments)
							if result.Error != nil {
								log.Printf("Error fetching comments for post ID %v: %v", post.PostID, result.Error)
								return nil, result.Error
							}
							log.Printf("Fetched %d comments for post ID %v", len(comments), post.PostID)
							return comments, nil
						}
						log.Printf("Failed to fetch comments: post not found or incorrect type")
						return nil, fmt.Errorf("could not retrieve comments for post")
					},
				},
				"user": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if post, ok := p.Source.(*models.Post); ok {
							// Возвращаем ID пользователя, связанного с постом
							return post.UserID, nil
						}
						log.Printf("Failed to fetch user ID: post not found or incorrect type")
						return nil, fmt.Errorf("could not retrieve user ID for post")
					},
				},
			},
		},
	)

	userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"firstName": &graphql.Field{
					Type: graphql.String,
				},
				"lastName": &graphql.Field{
					Type: graphql.String,
				},
				"posts": &graphql.Field{
					Type: graphql.NewList(postType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						user, ok := p.Source.(*models.User)
						if !ok || user == nil {
							log.Printf("Source conversion problem in posts field of userType")
							return nil, fmt.Errorf("source conversion problem")
						}
						var posts []models.Post
						// Используем db для выполнения запроса к базе данных
						err := db.Where("user_id = ?", user.UserID).Find(&posts).Error
						if err != nil {
							log.Printf("Error fetching posts for user ID %v: %v", user.UserID, err)
							return nil, err
						}
						log.Printf("Fetched %d posts for user ID %v", len(posts), user.UserID)
						return posts, nil
					},
				},
			},
		},
	)

	queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"posts": &graphql.Field{
					Type: graphql.NewList(postType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var posts []*models.Post
						log.Println("in queryType, before result := db.Find(&posts). OK")
						result := db.Find(&posts)
						log.Printf("in queryType, after result := db.Find(&posts). result = %#v\n", result)
						if result.Error != nil {
							log.Printf("Error fetching posts: %v", result.Error)
							return nil, result.Error
						}
						for _, post := range posts {
							log.Printf("Post: %#v", post)
						}
						log.Printf("Fetched %d posts", len(posts))
						return posts, nil
					},
				},
				"post": &graphql.Field{
					Type: postType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						var post models.Post
						id, ok := p.Args["id"].(int)
						if !ok {
							return nil, fmt.Errorf("id argument type must be int")
						}

						// Выполнение запроса к базе данных для поиска поста по ID
						result := db.First(&post, "post_id = ?", id)
						if result.Error != nil {
							if errors.Is(result.Error, gorm.ErrRecordNotFound) {
								return nil, fmt.Errorf("post with id '%v' not found", id)
							}
							// В случае других ошибок базы данных
							return nil, result.Error
						}

						// Возвращаем найденный пост
						return &post, nil
					},
				},
			},
		},
	)

	mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"updatePostCommentable": &graphql.Field{
				Type: postType, // Тип, который возвращается мутацией
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"commentable": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Boolean),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					postID := p.Args["id"].(int) // Приведение типа в соответствии с моделью
					commentable := p.Args["commentable"].(bool)

					// Найти пост по ID и обновить поле Commentable
					var post models.Post
					result := db.First(&post, postID)
					if result.Error != nil {
						if errors.Is(result.Error, gorm.ErrRecordNotFound) {
							return nil, fmt.Errorf("post with id '%v' not found", postID)
						}
						return nil, result.Error
					}

					// Обновление поля Commentable
					post.Commentable = commentable
					result = db.Save(&post)
					if result.Error != nil {
						return nil, result.Error
					}

					return &post, nil
				},
			},
		},
	})
	//log.Printf("userType: %#v\n", userType)
	//log.Printf("postType: %#v\n", postType)
	//log.Printf("commentType: %#v\n", commentType)
	//log.Printf("queryType: %#v\n", queryType)
}

// Рекурсивная функция для извлечения дочерних комментариев
func getNestedComments(db *gorm.DB, parentID uint) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := db.Where("parent_comment_id = ?", parentID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		children, err := getNestedComments(db, comment.CommentID)
		if err != nil {
			return nil, err
		}
		comment.ChildComments = children
	}

	return comments, nil
}
