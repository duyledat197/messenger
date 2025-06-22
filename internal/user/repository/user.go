package repository

import (
	"context"

	"openmyth/messgener/internal/user/entity"
	"openmyth/messgener/util/database"
)

// UserRepository is an interface that defines the contract for user repository operations.
type UserRepository interface {
	// Create creates a new user in the repository.
	Create(ctx context.Context, db database.Executor, data *entity.User) error

	// UpdateInfoByID updates the information for a user with the specified ID.
	UpdateInfoByID(ctx context.Context, db database.Executor, id string, data *entity.User) error

	RetrieveByID(ctx context.Context, db database.Executor, id string) (*entity.User, error)
	RetrieveByUserName(ctx context.Context, db database.Executor, userName string) (*entity.User, error)
	RetrieveByPhone(ctx context.Context, db database.Executor, phone string) (*entity.User, error)
}
