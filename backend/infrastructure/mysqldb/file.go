package mysqldb

import (
	"context"
	"database/sql"

	"github.com/kawabatas/go-fileuploader/domain/model"
	"github.com/kawabatas/go-fileuploader/domain/repository"
)

func NewFileMetadataRepository(db *sql.DB) *fileMetadataRepository {
	return &fileMetadataRepository{
		db: db,
	}
}

type fileMetadataRepository struct {
	db *sql.DB
}

var _ repository.FileMetadataRepository = (*fileMetadataRepository)(nil)

func (r *fileMetadataRepository) List(ctx context.Context, offset, limit int) ([]*model.File, error) {
	files := []*model.File{}
	query := "SELECT id, title FROM files ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		// 0件の場合の挙動、確認済み
		var file model.File
		if err := rows.Scan(&file.ID, &file.Title); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileMetadataRepository) Find(ctx context.Context, id string) (*model.File, error) {
	var file model.File
	var userID int
	query := "SELECT id, filename, title, size, mimetype, user_id, created_at FROM files WHERE id = ?"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&file.ID, &file.FileName, &file.Title, &file.Size, &file.MimeType, &userID, &file.CreatedAt); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(file.ID) == 0 {
		return nil, model.ErrNotFound
	}
	file.User = &model.User{
		ID: userID,
	}
	return &file, nil
}

func (r *fileMetadataRepository) Create(ctx context.Context, file *model.File) error {
	query := "INSERT INTO files (id, filename, title, size, mimetype, user_id) VALUES (?, ?, ?, ?, ?, ?)"
	if _, err := r.db.ExecContext(ctx, query, file.ID, file.FileName, file.Title, file.Size, file.MimeType, file.User.ID); err != nil {
		return err
	}
	return nil
}

func (r *fileMetadataRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM files WHERE id = ?"
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
