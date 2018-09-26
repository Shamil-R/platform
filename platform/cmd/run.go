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
	"github.com/vektah/gqlparser/ast"

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
		// dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		// 	"sa",
		// 	"p@sSw0rd",
		// 	"127.0.0.1",
		// 	1433,
		// 	"platform",
		// )

		// db, err := sqlx.Connect("mssql", dsn)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		conf := graph.Config{
			Resolvers: &platform.Resolver{},
		}

		router := chi.NewRouter()

		// router.Use(Middleware(db))

		router.Handle("/", handler.Playground("Platform", "/query"))
		router.Handle("/query",
			handler.GraphQL(
				graph.NewExecutableSchema(conf),
				handler.ResolverMiddleware(ResolverMiddleware()),
				handler.RequestMiddleware(RequestMiddleware()),
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

func ResolverMiddleware() graphql.FieldMiddleware {
	return func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		fmt.Println("ResolverMiddleware")

		// reqCtx := graphql.GetRequestContext(ctx)

		// for _, op := range reqCtx.Doc.Operations {
		// 	fmt.Println("op", op.Operation, op.Name)
		// }

		// if len(reqCtx.Errors) > 0 {
		// 	return nil, fmt.Errorf("error pizdec")
		// }

		resCtx := graphql.GetResolverContext(ctx)

		fmt.Println(resCtx.Object, "-", resCtx.Field.Name, "-", resCtx.Parent.Object)

		// fmt.Println("name", resCtx.Field.Name)

		// if strings.HasPrefix(resCtx.Field.Field.Name, "__") {
		// 	return next(ctx)
		// }

		// for _, op := range reqCtx.Doc.Operations {
		// 	fmt.Println(op.Operation, op.Name)
		// 	for _, sel := range op.SelectionSet {
		// 		fmt.Println("!", sel)
		// 		if f, ok := sel.(*ast.Field); ok {
		// 			fmt.Println("#", f.Name)
		// 		}
		// 	}
		// }

		// fieldSelections := resCtx.Field.Selections
		// var sels []string
		// for _, sel := range fieldSelections {
		// 	switch sel := sel.(type) {
		// 	case *ast.Field:
		// 		sels = append(sels, fmt.Sprintf("%s as %s in %s", sel.Name, sel.Alias, sel.ObjectDefinition.Name))
		// 	case *ast.InlineFragment:
		// 		sels = append(sels, fmt.Sprintf("inline fragment on %s", sel.TypeCondition))
		// 	case *ast.FragmentSpread:
		// 		// fragment := reqCtx.Doc.Fragments.ForName(sel.Name)
		// 		// sels = append(sels, fmt.Sprintf("named fragment %s on %s", sel.Name, fragment.TypeCondition))
		// 	}
		// }

		// fmt.Println(sels)

		res, err := next(ctx)

		// fmt.Println("res", res)

		return res, err
	}
}

func RequestMiddleware() graphql.RequestMiddleware {
	return func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		reqCtx := graphql.GetRequestContext(ctx)

		fmt.Println("RequestMiddleware")

		for _, op := range reqCtx.Doc.Operations {
			fmt.Println("op", op.Operation, op.Name)
			tr(op.SelectionSet)
		}

		reqCtx.Error(ctx, fmt.Errorf("ERROR"))

		return next(ctx)
	}
}

func tr(set ast.SelectionSet) {
	if set != nil {
		for _, sel := range set {
			if f, ok := sel.(*ast.Field); ok {
				fmt.Println("sel", f.Name)
				tr(f.SelectionSet)
			}
		}
	}
}
