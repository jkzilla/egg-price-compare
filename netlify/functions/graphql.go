package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/jkzilla/egg-price-compare/graph"
)

var graphqlHandler *httpadapter.HandlerAdapter

func init() {
	resolver := graph.NewResolver()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	graphqlHandler = httpadapter.New(srv)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return graphqlHandler.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
