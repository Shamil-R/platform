package cmd

import (
	"context"
	"fmt"
	"gitlab/nefco/platform"
	"gitlab/nefco/platform/server/graph"
	"gitlab/nefco/platform/service"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"

	"github.com/99designs/gqlgen/graphql"
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
				handler.ResolverMiddleware(ResolverMiddleware(service.NewRoleService())),
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

func ResolverMiddleware(roleService service.RoleService) graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		// var sels []string

		// reqCtx := graphql.GetRequestContext(ctx)
		resCtx := graphql.GetResolverContext(ctx)

		roleService.CheckRole(resCtx)

		// fieldSelections := resCtx.Field.Selections

		// for _, sel := range fieldSelections {
		// 	switch sel := sel.(type) {
		// 	case *ast.Field:
		// 		sels = append(sels, fmt.Sprintf("%s as %s in %s", sel.Name, sel.Alias, sel.ObjectDefinition.Name))
		// 	case *ast.InlineFragment:
		// 		sels = append(sels, fmt.Sprintf("inline fragment on %s", sel.TypeCondition))
		// 	case *ast.FragmentSpread:
		// 		fragment := reqCtx.Doc.Fragments.ForName(sel.Name)
		// 		sels = append(sels, fmt.Sprintf("named fragment %s on %s", sel.Name, fragment.TypeCondition))
		// 	}
		// }

		// fmt.Println(sels)

		return next(ctx)
	}
}
