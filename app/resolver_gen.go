//go:generate gorunpkg github.com/99designs/gqlgen

package app

import (
	context "context"
	graph "gitlab/nefco/platform/app/graph"
	model "gitlab/nefco/platform/app/model"
)

type Resolver struct{}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) User() graph.UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateMaterial(ctx context.Context, data model.MaterialCreateInput) (model.Material, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateMaterial(ctx context.Context, data model.MaterialUpdateInput, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteMaterial(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpsertMaterial(ctx context.Context, where model.MaterialWhereUniqueInput, create model.MaterialCreateInput, update model.MaterialUpdateInput) (*model.Material, error) {
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
func (r *mutationResolver) UpsertUser(ctx context.Context, where model.UserWhereUniqueInput, create model.UserCreateInput, update model.UserUpdateInput) (*model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreatePerson(ctx context.Context, data model.PersonCreateInput) (model.Person, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdatePerson(ctx context.Context, data model.PersonUpdateInput, where model.PersonWhereUniqueInput) (*model.Person, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeletePerson(ctx context.Context, where model.PersonWhereUniqueInput) (*model.Person, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpsertPerson(ctx context.Context, where model.PersonWhereUniqueInput, create model.PersonCreateInput, update model.PersonUpdateInput) (*model.Person, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Material(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	panic("not implemented")
}
func (r *queryResolver) Materials(ctx context.Context, where *model.MaterialWhereInput, skip *int, first *int, last *int) ([]*model.Material, error) {
	panic("not implemented")
}
func (r *queryResolver) User(ctx context.Context, where model.UserWhereUniqueInput) (*model.User, error) {
	panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context, where *model.UserWhereInput, skip *int, first *int, last *int) ([]*model.User, error) {
	panic("not implemented")
}
func (r *queryResolver) Person(ctx context.Context, where model.PersonWhereUniqueInput) (*model.Person, error) {
	panic("not implemented")
}
func (r *queryResolver) People(ctx context.Context, where *model.PersonWhereInput, skip *int, first *int, last *int) ([]*model.Person, error) {
	panic("not implemented")
}

type userResolver struct{ *Resolver }

func (r *userResolver) Materials(ctx context.Context, obj *model.User, where *model.MaterialWhereInput) ([]model.Material, error) {
	panic("not implemented")
}
