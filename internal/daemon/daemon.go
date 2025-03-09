package daemon

import (
	"context"

	"github.com/cepmap/otus-system-monitoring/internal/storage"
)

type iStorage struct {
	loadAvgStorage  storage.Storage
	cpuStorage      storage.Storage
	diskLoadStorage storage.Storage
	diskInfoStorage storage.Storage
}
type Server struct {
	storage *iStorage
	ctx     context.Context
}
