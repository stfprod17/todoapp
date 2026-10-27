package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/stfprod17/todoapp/internal/core/domain"
	core_errors "github.com/stfprod17/todoapp/internal/core/errors"
	core_postgres_pool "github.com/stfprod17/todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query :=
		`
    SELECT id, version, full_name, phone_number, email, password_hash
    FROM todoapp.users
    WHERE LOWER(email)=LOWER($1);
    
`

	row := r.pool.QueryRow(ctx, query, email)

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with email=%q: %w", email, core_errors.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("scan user by email: %w", err)
	}

	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
		userModel.Email,
		userModel.PasswordHash,
	), nil
}
