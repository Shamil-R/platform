package service

import (
	"context"
	model "gitlab/nefco/platform/server/model"
	"log"

	"github.com/jmoiron/sqlx"
)

type service struct {
}

func NewService() {
}

func (s *service) CreatePost(ctx context.Context, data model.PostCreateInput) (*model.Post, error) {
	tx, ok := ctx.Value("ctx_tx").(*sqlx.Tx)
	if !ok {
		log.Fatal("ctx_tx failed")
	}
	return nil, nil
}
