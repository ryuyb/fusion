package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/domain/entity"
	"github.com/ryuyb/fusion/internal/domain/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/user"
	"github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type userRepository struct {
	client *database.Client
	logger *zap.Logger
}

func NewUserRepository(client *database.Client, logger *zap.Logger) repository.UserRepository {
	return &userRepository{
		client: client,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	created, err := r.client.User.
		Create().
		SetUsername(u.Username).
		SetPassword(u.Password).
		SetEmail(u.Email).
		SetStatus(user.Status(u.Status)).
		Save(ctx)
	if err != nil {
		r.logger.Error("failed to create user", zap.Error(err), zap.String("username", u.Username), zap.String("email", u.Email))
		return nil, errors.ConvertDatabaseError(err, "User")
	}
	return r.toEntity(created), nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("User").WithDetail("id", id)
		}
		r.logger.Error("failed to find user by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors.ConvertDatabaseError(err, "User")
	}
	return r.toEntity(u), nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("User").WithDetail("username", username)
		}
		r.logger.Error("failed to find user by username", zap.Error(err), zap.String("username", username))
		return nil, errors.ConvertDatabaseError(err, "User")
	}
	return r.toEntity(u), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("User").WithDetail("email", email)
		}
		r.logger.Error("failed to find user by email", zap.Error(err), zap.String("email", email))
		return nil, errors.ConvertDatabaseError(err, "User")
	}
	return r.toEntity(u), nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("User").WithDetail("id", id)
		}
		r.logger.Error("failed to delete user", zap.Error(err), zap.Int64("id", id))
		return errors.ConvertDatabaseError(err, "User")
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, u *entity.User) (*entity.User, error) {
	updated, err := r.client.User.
		UpdateOneID(u.ID).
		SetUsername(u.Username).
		SetPassword(u.Password).
		SetEmail(u.Email).
		SetStatus(user.Status(u.Status)).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("User").WithDetail("id", u.ID)
		}
		r.logger.Error("failed to update user", zap.Error(err), zap.Int64("id", u.ID))
		return nil, errors.ConvertDatabaseError(err, "User")
	}
	return r.toEntity(updated), nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int, error) {
	total, err := r.client.User.Query().Count(ctx)
	if err != nil {
		r.logger.Error("failed to count users", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	users, err := r.client.User.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list users", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	entities := make([]*entity.User, len(users))
	for i, u := range users {
		entities[i] = r.toEntity(u)
	}

	return entities, total, nil
}

func (r *userRepository) toEntity(u *ent.User) *entity.User {
	return &entity.User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Status:    entity.UserStatus(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeleteAt:  u.DeleteAt,
	}
}
