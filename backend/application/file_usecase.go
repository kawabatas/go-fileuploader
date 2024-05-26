package application

import (
	"context"
	"io"
	"path/filepath"

	"github.com/kawabatas/go-fileuploader/domain/model"
	"github.com/kawabatas/go-fileuploader/domain/repository"
)

func NewFileUseCase(
	metadataRepo repository.FileMetadataRepository,
	bucketRepository repository.FileBucketRepository,
	bucketBaseURL string,
) *FileUseCase {
	return &FileUseCase{
		metadataRepo:  metadataRepo,
		bucketRepo:    bucketRepository,
		bucketBaseURL: bucketBaseURL,
	}
}

type FileUseCase struct {
	metadataRepo  repository.FileMetadataRepository
	bucketRepo    repository.FileBucketRepository
	bucketBaseURL string
}

func (uc *FileUseCase) GetList(ctx context.Context, offset, limit int) ([]*model.File, error) {
	files, err := uc.metadataRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (uc *FileUseCase) GetDetail(ctx context.Context, id string) (*model.File, error) {
	file, err := uc.metadataRepo.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	file.SetBaseURL(uc.bucketBaseURL)
	return file, nil
}

func (uc *FileUseCase) Post(
	ctx context.Context,
	file io.ReadCloser,
	fileName, title string,
	fileSize int,
	mimeType string,
) error {
	f, err := model.NewFile()
	if err != nil {
		return err
	}

	ext := filepath.Ext(fileName)
	key := f.ID + ext
	// ファイルオブジェクトをアップロード
	if err := uc.bucketRepo.Upload(ctx, key, file); err != nil {
		return err
	}
	// NOTE: オブジェクトのアップロードと、メタデータの保存の間で、エラーが発生する場合も考えられるが、許容する

	// メタデータを保存
	f.FileName = fileName
	f.Title = title
	f.Size = fileSize
	f.MimeType = mimeType
	if err := uc.metadataRepo.Create(ctx, f); err != nil {
		return err
	}
	return nil
}

func (uc *FileUseCase) Delete(ctx context.Context, id string) error {
	// メタデータを削除
	if err := uc.metadataRepo.Delete(ctx, id); err != nil {
		return err
	}
	// NOTE: メタデータの削除と、オブジェクトの削除の間で、エラーが発生する場合も考えられるが、許容する

	// ファイルオブジェクトを削除
	key := id + ".*"
	if err := uc.bucketRepo.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
