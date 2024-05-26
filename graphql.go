package main

import (
	"OZON/models"
	"fmt"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"log"
)

var (
	db          *gorm.DB
	userType    *graphql.Object
	postType    *graphql.Object
	commentType *graphql.Object
	queryType   *graphql.Object
)

func initGraphQLTypes(db *gorm.DB) {
	commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"content": &graphql.Field{
					Type: graphql.String,
				},
				"createdTime": &graphql.Field{
					Type: graphql.DateTime,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if comment, ok := p.Source.(*models.Comment); ok && comment != nil {
							log.Printf("Retrieving created time for comment ID: %v", comment.ID)
							return comment.CreatedAt, nil
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
			},
		},
	)

	postType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Post",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
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
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if post, ok := p.Source.(*models.Post); ok {
							var comments []models.Comment
							result := db.Where("post_id = ?", post.ID).Find(&comments)
							if result.Error != nil {
								log.Printf("Error fetching comments for post ID %v: %v", post.ID, result.Error)
								return nil, result.Error
							}
							log.Printf("Fetched %d comments for post ID %v", len(comments), post.ID)
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
						err := db.Where("user_id = ?", user.ID).Find(&posts).Error
						if err != nil {
							log.Printf("Error fetching posts for user ID %v: %v", user.ID, err)
							return nil, err
						}
						log.Printf("Fetched %d posts for user ID %v", len(posts), user.ID)
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
						var posts []models.Post
						// Выполняем запрос к базе данных, чтобы найти все посты
						log.Println("in queryType, before result := db.Find(&posts). OK")
						result := db.Find(&posts)
						log.Printf("in queryType, after result := db.Find(&posts). result = %#v\n", result)
						if result.Error != nil {
							log.Printf("Error fetching posts: %v", result.Error)
							return nil, result.Error
						}
						log.Printf("Fetched %d posts", len(posts))
						return posts, nil
					},
				},
			},
		},
	)

	//log.Printf("userType: %#v\n", userType)
	//log.Printf("postType: %#v\n", postType)
	//log.Printf("commentType: %#v\n", commentType)
	//log.Printf("queryType: %#v\n", queryType)
}
