package service

import (
	"context"
	"demo/internal/common"
	historyModel "demo/internal/domain/history/model"
	"demo/internal/domain/sk/model"
	"demo/internal/domain/sk/storage"
	storageModel "demo/internal/domain/sk/storage/model"
	"demo/internal/util"
	"errors"
)

type Sk struct {
	skStorage      *storage.Sk
	historyService History
}

func NewSk(skStorage *storage.Sk, history History) *Sk {
	return &Sk{
		skStorage:      skStorage,
		historyService: history,
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

	err = s.historyService.SaveAction(ctx, historyModel.HistoryAction{
		Id:    sk.Id,
		Type:  historyModel.Add,
		Table: historyModel.Sk,
	})
	if err != nil {
		return common.Wrap(err, "SaveAction")
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

	err = s.historyService.SaveAction(ctx, historyModel.HistoryAction{
		Id:    sk.Id,
		Type:  historyModel.Edit,
		Table: historyModel.Sk,
	})
	if err != nil {
		return common.Wrap(err, "SaveAction")
	}

	return nil
}
