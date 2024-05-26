package main

import (
	"OZON/pkg/database"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
)

func main() {
	var err error
	db, err = database.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("Database connection established: %#v\n", db)

	fmt.Println("DB = ", db)
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	// Инициализация GraphQL типов с передачей db
	initGraphQLTypes(db)

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
			// Mutation: mutationType,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	// HTTP-сервер для обработки запросов GraphQL
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)

	// Запустите HTTP-сервер
	log.Println("Starting server on :8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
