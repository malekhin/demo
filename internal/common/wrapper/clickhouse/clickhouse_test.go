package clickhouse

import (
	"context"
	"demo/internal/common/config"
	mocks2 "demo/internal/common/metrics/mocks"
	"demo/pkg/clickhouse"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestClickhouse(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)

	pkgCllick := clickhouse.New(cfg, zap.NewNop())
	mericsMock := mocks2.NewIMetrics(t)

	mericsMock.On("RecordExternalCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	click := New(pkgCllick, mericsMock)

	err = click.Exec(context.Background(), `select true`)
	require.NoError(t, err)

	err = click.Exec(context.Background(), `syntax err`)
	require.Error(t, err)
}
