package service

import (
	"context"
	"demo/internal/domain/history/model"
)

type History interface {
	SaveAction(ctx context.Context, historyAction model.HistoryAction) error
}
