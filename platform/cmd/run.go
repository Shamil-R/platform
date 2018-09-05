package cmd

import (
	"context"
	"errors"
	"gitlab/nefco/platform"
	"gitlab/nefco/platform/server/graph"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/jmoiron/sqlx"

	"github.com/99designs/gqlgen/handler"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `run`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := graph.Config{
			Resolvers: &platform.Resolver{},
		}
		http.Handle("/", handler.Playground("Todo", "/query"))
		http.Handle("/query", handler.GraphQL(
			graph.NewExecutableSchema(conf),
			handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
				// send this panic somewhere
				log.Print(err)
				debug.PrintStack()
				return errors.New("user message on panic")
			}),
		))
		log.Fatal(http.ListenAndServe(":8081", nil))
	},
}

func Middleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			ctx := context.WithValue(r.Context(), "ctx_tx", tx)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
