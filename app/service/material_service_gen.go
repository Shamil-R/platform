// Code generated by platform, DO NOT EDIT.

package service

import (
	"context"

	model "gitlab/nefco/platform/app/model"

	mssql "gitlab/nefco/platform/codegen/generate/service/mssql"
)

type MaterialMutationService interface {
	CreateMaterial(ctx context.Context, data model.MaterialCreateInput) (model.Material, error)
	UpdateMaterial(ctx context.Context, data model.MaterialUpdateInput, where model.MaterialWhereUniqueInput) (*model.Material, error)
	DeleteMaterial(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error)
	UpsertMaterial(ctx context.Context, where model.MaterialWhereUniqueInput, create model.MaterialCreateInput, update model.MaterialUpdateInput) (*model.Material, error)
}

type materialMutationService struct{}

func NewMaterialMutationService() *materialMutationService {
	return &materialMutationService{}
}

func (s *materialMutationService) CreateMaterial(ctx context.Context, data model.MaterialCreateInput) (model.Material, error) {
	res := model.Material{}

	if err := mssql.Create(ctx, &res); err != nil {
		return res, err
	}

	return res, nil
}

func (s *materialMutationService) UpdateMaterial(ctx context.Context, data model.MaterialUpdateInput, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	res := &model.Material{}

	if err := mssql.Update(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *materialMutationService) DeleteMaterial(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	res := &model.Material{}

	if err := mssql.Delete(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
	
}

func (s *materialMutationService) UpsertMaterial(ctx context.Context, where model.MaterialWhereUniqueInput, create model.MaterialCreateInput, update model.MaterialUpdateInput) (*model.Material, error) {
	panic("not implemented")
}

type MaterialQueryService interface {
	Material(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error)
	Materials(ctx context.Context, where *model.MaterialWhereInput, skip *int, first *int, last *int) ([]*model.Material, error)
}

type materialQueryService struct{}

func NewMaterialQueryService() *materialQueryService {
	return &materialQueryService{}
}

func (s *materialQueryService) Material(ctx context.Context, where model.MaterialWhereUniqueInput) (*model.Material, error) {
	res := &model.Material{}

	if err := mssql.Item(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *materialQueryService) Materials(ctx context.Context, where *model.MaterialWhereInput, skip *int, first *int, last *int) ([]*model.Material, error) {
	res := []*model.Material{}

	if err := mssql.Collection(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

type MaterialRelationService interface {
	Author(ctx context.Context, obj *model.Material, where *model.UserWhereInput) (model.User, error)
}

type materialRelationService struct{}

func NewMaterialRelationService() *materialRelationService {
	return &materialRelationService{}
}

func (s *materialRelationService) Author(ctx context.Context, obj *model.Material, where *model.UserWhereInput) (model.User, error) {
	res := model.User{}

	if err := mssql.Relation(ctx, obj.ID, &res); err != nil {
		return res, err
	}

	return res, nil
}
