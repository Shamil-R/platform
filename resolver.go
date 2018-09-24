package platform

import (
	context "context"
	graph "gitlab/nefco/platform/server/graph"
	model "gitlab/nefco/platform/server/model"
)

type Resolver struct{}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data model.UserCreateInput) (model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateUser(ctx context.Context, data model.UserUpdateInput, where model.UserWhereUniqueInput) (*model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteUser(ctx context.Context, where model.UserWhereUniqueInput) (*model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateRoleObject(ctx context.Context, data model.RoleObjectCreateInput) (model.RoleObject, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateRoleObject(ctx context.Context, data model.RoleObjectUpdateInput, where model.RoleObjectWhereUniqueInput) (*model.RoleObject, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteRoleObject(ctx context.Context, where model.RoleObjectWhereUniqueInput) (*model.RoleObject, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateMaterial(ctx context.Context, data model.MaterialCreateInput) (model.Material, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateMaterial(ctx context.Context, data model.MaterialUpdateInput, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteMaterial(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, where model.UserWhereUniqueInput) (*model.User, error) {
	panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context, where *model.UserWhereInput) ([]*model.User, error) {
	panic("not implemented")
}
func (r *queryResolver) RoleObject(ctx context.Context, where model.RoleObjectWhereUniqueInput) (*model.RoleObject, error) {
	panic("not implemented")
}
func (r *queryResolver) RoleObjects(ctx context.Context, where *model.RoleObjectWhereInput) ([]*model.RoleObject, error) {
	panic("not implemented")
}
func (r *queryResolver) Material(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	panic("not implemented")
}
func (r *queryResolver) Materials(ctx context.Context, where *model.MaterialWhereInput) ([]*model.Material, error) {
	panic("not implemented")
}
