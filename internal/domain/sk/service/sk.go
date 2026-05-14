package service

import (
	"context"
	"demo/internal/common"
	"demo/internal/domain/sk/model"
	"demo/internal/domain/sk/storage"
	storageModel "demo/internal/domain/sk/storage/model"
	"demo/internal/util"
	"errors"
)

type Sk struct {
	skStorage *storage.Sk
}

func NewSk(skStorage *storage.Sk) *Sk {
	return &Sk{
		skStorage: skStorage,
	}
}
func (s *Sk) List(ctx context.Context, filter model.SkFilter) ([]model.SkItem, int, error) {
	list, count, err := s.skStorage.List(ctx, filter)
	if err != nil {
		return nil, 0, common.Wrap(err, "skStorage.List")
	}

	res := make([]model.SkItem, 0, len(list))
	for _, sk := range list {
		res = append(res, model.SkItem{
			Id:       sk.Id,
			Name:     sk.Name,
			IsActive: sk.IsActive,
		})
	}

	return res, count, nil
}

func (s *Sk) Add(ctx context.Context, request model.SkAdd) error {
	sk := storageModel.Sk{
		Id:       request.Id,
		Name:     request.Name,
		IsActive: request.IsActive,
	}

	isDuplicate, err := s.skStorage.IsDuplicate(ctx, sk)
	if err != nil {
		return common.Wrap(err, "Add")
	}
	if isDuplicate {
		return util.NewBadRequest(errors.New("duplicate"))
	}

	_, err = s.skStorage.Add(ctx, sk)
	if err != nil {
		return common.Wrap(err, "Add")
	}

	return nil
}

func (s *Sk) Edit(ctx context.Context, id int, request model.SkEdit) error {
	sk := storageModel.Sk{
		Id:       id,
		Name:     request.Name,
		IsActive: request.IsActive,
	}

	_, err := s.skStorage.Edit(ctx, sk)
	if err != nil {
		return common.Wrap(err, "Edit")
	}

	return nil
}
