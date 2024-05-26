package localfs

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/kawabatas/go-fileuploader/domain/repository"
)

func NewFileBucketRepository(basePath string) *fileBucketRepository {
	return &fileBucketRepository{
		basePath: basePath,
	}
}

type fileBucketRepository struct {
	basePath string
}

var _ repository.FileBucketRepository = (*fileBucketRepository)(nil)

func (r *fileBucketRepository) Upload(ctx context.Context, key string, src io.ReadCloser) error {
	dst, err := os.Create(filepath.Join(r.basePath, key))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (r *fileBucketRepository) Delete(ctx context.Context, key string) error {
	files, err := filepath.Glob(filepath.Join(r.basePath, key))
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return nil
		}
	}
	return nil
}
