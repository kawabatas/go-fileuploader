package repository

import (
	"context"

	"github.com/kawabatas/go-fileuploader/domain/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}
