package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jkzilla/egg-price-compare/graph"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/", playground.Handler("Egg Price Comparison", "/graphql"))
	http.Handle("/graphql", c.Handler(srv))

	log.Printf("🥚 Egg Price Comparison API")
	log.Printf("GraphQL playground: http://localhost:%s/", port)
	log.Printf("GraphQL endpoint: http://localhost:%s/graphql", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
