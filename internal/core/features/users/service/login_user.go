package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/stfprod17/todoapp/internal/core/domain"
	core_errors "github.com/stfprod17/todoapp/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) LoginUser(
	ctx context.Context,
	email string,
	password string,
) (domain.User, error) {
	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return domain.User{}, fmt.Errorf("invalid email or password: %w", core_errors.ErrInvalidArgument)
		}

		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, fmt.Errorf("invalid email or password: %w", core_errors.ErrInvalidArgument)
	}

	return user, nil
}
