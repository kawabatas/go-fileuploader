package application

import (
	"context"
	"errors"
	"io"
	"path/filepath"

	"github.com/kawabatas/go-fileuploader/domain/model"
	"github.com/kawabatas/go-fileuploader/domain/repository"
)

func NewFileUseCase(
	metadataRepo repository.FileMetadataRepository,
	bucketRepository repository.FileBucketRepository,
	userRepository repository.UserRepository,
	bucketBaseURL string,
) *FileUseCase {
	return &FileUseCase{
		metadataRepo:   metadataRepo,
		bucketRepo:     bucketRepository,
		userRepository: userRepository,
		bucketBaseURL:  bucketBaseURL,
	}
}

type FileUseCase struct {
	metadataRepo   repository.FileMetadataRepository
	bucketRepo     repository.FileBucketRepository
	userRepository repository.UserRepository
	bucketBaseURL  string
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

	user, err := uc.userRepository.FindByID(ctx, file.User.ID)
	if err != nil {
		return nil, err
	}
	file.SetUser(user)
	return file, nil
}

func (uc *FileUseCase) Post(
	ctx context.Context,
	file io.ReadCloser,
	fileName, title string,
	fileSize int,
	mimeType string,
	userEmail string,
) error {
	var user *model.User
	var err error
	user, err = uc.userRepository.FindByEmail(ctx, userEmail)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return err
		}
	}
	if user == nil {
		user, err = uc.userRepository.Create(ctx, &model.User{Email: userEmail})
		if err != nil {
			return err
		}
	}

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
	f.SetUser(user)
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
