package mysqldb

import (
	"context"
	"database/sql"

	"github.com/kawabatas/go-fileuploader/domain/model"
	"github.com/kawabatas/go-fileuploader/domain/repository"
)

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *sql.DB
}

var _ repository.UserRepository = (*userRepository)(nil)

func (r *userRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	query := "SELECT id, email FROM users WHERE id = ?"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, model.ErrNotFound
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := "SELECT id, email FROM users WHERE email = ?"
	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, model.ErrNotFound
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := "INSERT INTO users (email) VALUES (?)"
	if _, err := r.db.ExecContext(ctx, query, user.Email); err != nil {
		return nil, err
	}
	return r.FindByEmail(ctx, user.Email)
}
