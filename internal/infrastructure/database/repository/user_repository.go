package repository

import (
	"context"

	"github.com/ryuyb/fusion/internal/core/domain"
	"github.com/ryuyb/fusion/internal/core/port/repository"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
	"github.com/ryuyb/fusion/internal/infrastructure/database/ent/user"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
	"go.uber.org/zap"
)

type userRepository struct {
	client *ent.Client
	logger *zap.Logger
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	created, err := r.client.User.
		Create().
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		r.logger.Error("failed to create user",
			zap.Error(err),
			zap.String("username", user.Username),
			zap.String("email", user.Email),
		)
		return nil, errors2.ConvertDatabaseError(err, "User")
	}
	return r.toDomain(created), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	updated, err := r.client.User.
		UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("User").WithDetail("id", user.ID)
		}
		r.logger.Error("failed to update user", zap.Error(err), zap.String("username", user.Username), zap.String("email", user.Email))
		return nil, errors2.ConvertDatabaseError(err, "User")
	}
	return r.toDomain(updated), nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors2.NotFound("User").WithDetail("id", id)
		}
		r.logger.Error("failed to delete user", zap.Error(err), zap.Int64("id", id))
		return errors2.ConvertDatabaseError(err, "User")
	}
	return nil
}

func (r *userRepository) FindById(ctx context.Context, id int64) (*domain.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("User").WithDetail("id", id)
		}
		r.logger.Error("failed to find user by id", zap.Error(err), zap.Int64("id", id))
		return nil, errors2.ConvertDatabaseError(err, "User")
	}
	return r.toDomain(u), nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("User").WithDetail("username", username)
		}
		r.logger.Error("failed to find user by username", zap.Error(err), zap.String("username", username))
		return nil, errors2.ConvertDatabaseError(err, "User")
	}
	return r.toDomain(u), nil
}

func (r *userRepository) ExistByUsername(ctx context.Context, username string) (bool, error) {
	exist, err := r.client.User.
		Query().
		Where(user.UsernameEQ(username)).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to find user by username", zap.Error(err), zap.String("username", username))
		return false, errors2.ConvertDatabaseError(err, "User")
	}
	return exist, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors2.NotFound("User").WithDetail("email", email)
		}
		r.logger.Error("failed to find user by email", zap.Error(err), zap.String("email", email))
		return nil, errors2.ConvertDatabaseError(err, "User")
	}
	return r.toDomain(u), nil
}

func (r *userRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	exist, err := r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)
	if err != nil {
		r.logger.Error("failed to find user by email", zap.Error(err), zap.String("email", email))
		return false, errors2.ConvertDatabaseError(err, "User")
	}
	return exist, nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int, error) {
	total, err := r.client.User.Query().Count(ctx)
	if err != nil {
		r.logger.Error("failed to count users", zap.Error(err))
		return nil, 0, errors2.DatabaseError(err)
	}
	users, err := r.client.User.
		Query().
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(user.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list users", zap.Error(err))
		return nil, 0, errors2.DatabaseError(err)
	}
	entities := make([]*domain.User, len(users))
	for i, u := range users {
		entities[i] = r.toDomain(u)
	}
	return entities, total, nil
}

func (r *userRepository) toDomain(user *ent.User) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserRepository(client *ent.Client, logger *zap.Logger) repository.UserRepository {
	return &userRepository{
		client: client,
		logger: logger,
	}
}
