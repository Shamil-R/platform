package cmd

import (
	"context"
	"fmt"
	"gitlab/nefco/platform"
	"gitlab/nefco/platform/server/graph"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"

	"github.com/99designs/gqlgen/handler"
	_ "github.com/denisenkom/go-mssqldb"
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
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			"sa",
			"p@sSw0rd",
			"127.0.0.1",
			1433,
			"platform",
		)

		db, err := sqlx.Connect("mssql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		conf := graph.Config{
			Resolvers: &platform.Resolver{},
		}

		router := chi.NewRouter()

		router.Use(Middleware(db))

		router.Handle("/", handler.Playground("Platform", "/query"))
		router.Handle("/query",
			handler.GraphQL(
				graph.NewExecutableSchema(conf),
				// handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
				// 	// send this panic somewhere
				// 	log.Print(err)
				// 	debug.PrintStack()
				// 	return errors.New("user message on panic")
				// }),
			))
		if err := http.ListenAndServe(":8080", router); err != nil {
			panic(err)
		}
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
