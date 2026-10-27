package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/stfprod17/todoapp/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number, email, password_hash)
    VALUES ($1, $2, $3, $4)
    RETURNING id, version, full_name, phone_number, email, password_hash;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		user.Email,
		user.PasswordHash,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.Email,
		&userModel.PasswordHash,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
		userModel.Email,
		userModel.PasswordHash,
	)

	return userDomain, nil
}
