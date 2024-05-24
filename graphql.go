package main

import (
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(
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
					// ЗАГЛУШКА. Тут логика для получения постов
					return nil, nil
				},
			},
		},
	},
)

var postType = graphql.NewObject(
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
			},
		},
	},
)

var commentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// логика для получения пользователя по user_id (заглушка)
					return nil, nil
				},
			},
		},
	},
)
