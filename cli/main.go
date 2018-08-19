package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/handler"

	"gitlab/nefco/platform"
	"gitlab/nefco/platform/graph/generated"
)

func main() {
	conf := generated.Config{
		Resolvers: &platform.Resolver{},
	}

	http.Handle("/", handler.Playground("Playground", "/query"))
	http.Handle("/query", handler.GraphQL(
		generated.NewExecutableSchema(conf),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			// send this panic somewhere
			log.Print(err)
			debug.PrintStack()
			return errors.New("user message on panic")
		}),
	))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
