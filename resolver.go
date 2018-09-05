// go:generate gorunpkg github.com/99designs/gqlgen

package platform

import (
	context "context"
	"fmt"
	graph "gitlab/nefco/platform/server/graph"
	model "gitlab/nefco/platform/server/model"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type Resolver struct {
	db *sqlx.DB
}

func NewResolver() *Resolver {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		"sa",
		"p@sSw0rd",
		"127.0.0.1",
		"1433",
		"platform",
	)
	db, err := sqlx.Connect("mssql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return &Resolver{
		db: db,
	}
}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, data model.PostCreateInput) (model.Post, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdatePost(ctx context.Context, data model.PostUpdateInput, where model.PostWhereUniqueInput) (*model.Post, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeletePost(ctx context.Context, where model.PostWhereUniqueInput) (*model.Post, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateUser(ctx context.Context, data model.UserCreateInput) (model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateUser(ctx context.Context, data model.UserUpdateInput, where model.UserWhereUniqueInput) (*model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteUser(ctx context.Context, where model.UserWhereUniqueInput) (*model.User, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Post(ctx context.Context, where model.PostWhereUniqueInput) (*model.Post, error) {
	panic("not implemented")
}
func (r *queryResolver) Posts(ctx context.Context, where *model.PostWhereInput) ([]*model.Post, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, where model.UserWhereUniqueInput) (*model.User, error) {
	user := &model.User{
		ID:   "1",
		Name: "Test",
	}
	return user, nil
	// panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context, where *model.UserWhereInput) ([]*model.User, error) {
	panic("not implemented")
}
