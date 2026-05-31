package worker_pool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWorkerPool(t *testing.T) {
	wp, err := New(-1, 5)
	require.Error(t, err)

	wp, err = New(5, 5)
	require.NoError(t, err)

	err = wp.AddTask(nil)
	require.Error(t, err)

	for i := 0; i < 6; i++ {
		err = wp.AddTask(func() {
			time.Sleep(time.Duration(i) * time.Second)
		})

		if i < 5 {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	err = wp.Close()
	require.NoError(t, err)
}
