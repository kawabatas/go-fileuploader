package repository

import (
	"context"
	"io"

	"github.com/kawabatas/go-fileuploader/domain/model"
)

type FileMetadataRepository interface {
	List(ctx context.Context, offset, limit int) ([]*model.File, error)
	Find(ctx context.Context, id string) (*model.File, error)
	Create(ctx context.Context, file *model.File) error
	Delete(ctx context.Context, id string) error
}

type FileBucketRepository interface {
	Upload(ctx context.Context, key string, src io.ReadCloser) error
	Delete(ctx context.Context, key string) error
}
