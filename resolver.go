// go:generate gorunpkg github.com/99designs/gqlgen

package platform

import (
	context "context"
	graph "gitlab/nefco/platform/server/graph"
	model "gitlab/nefco/platform/server/model"
)

type Resolver struct {
}

func NewResolver() *Resolver {
	return &Resolver{}
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

	// fmt.Println(rctx.Object)
	// fmt.Println(rctx.Field.Name)
	// for _, arg := range rctx.Field.Arguments {
	// 	fmt.Println(arg.Name)
	// }
	// fmt.Println(rctx.Field.Definition.Name)
	// for k, arg := range rctx.Args {
	// 	fmt.Println(k, arg)
	// }
	// rctx.Field.SelectionSet[0]
	// for _, s := range rctx.Field.Selections {
	// 	fmt.Println(s.GetPosition())
	// }

	// fmt.Println(rctx.Field.ObjectDefinition.Name)
	// for _, field := range rctx.Field.ObjectDefinition.Fields {
	// 	fmt.Println(field.Name, field.Type)
	// }

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
