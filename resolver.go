//go:generate gorunpkg github.com/99designs/gqlgen

package platform

import (
	context "context"
	"fmt"
	generated "gitlab/nefco/platform/graph/generated"
	models "gitlab/nefco/platform/models"

	"github.com/99designs/gqlgen/graphql"
)

type Resolver struct{}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context, where *models.UserWhereInput) ([]*models.User, error) {
	users := []*models.User{
		&models.User{"1", "Den"},
	}
	return users, nil
}

func (r *queryResolver) User(ctx context.Context, where models.UserWhereUniqueInput) (*models.User, error) {
	user := &models.User{"2", "Rent"}
	return user, nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data models.UserCreateInput) (models.User, error) {
	fmt.Println(graphql.GetResolverContext(ctx).Field.Selections)
	user := models.User{}
	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, data models.UserUpdateInput, where models.UserWhereUniqueInput) (*models.User, error) {
	return nil, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, where models.UserWhereUniqueInput) (*models.User, error) {
	return nil, nil
}
