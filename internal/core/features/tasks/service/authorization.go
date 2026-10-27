package tasks_service

import (
	"context"
	"fmt"

	core_auth "github.com/stfprod17/todoapp/internal/core/auth"
	"github.com/stfprod17/todoapp/internal/core/domain"
	core_errors "github.com/stfprod17/todoapp/internal/core/errors"
)

func currentUserID(ctx context.Context) (int, error) {
	userID, ok := core_auth.UserIDFromContext(ctx)
	if !ok {
		return 0, fmt.Errorf("user id is missing in context: %w", core_errors.ErrInvalidArgument)
	}

	return userID, nil
}

func ensureTaskOwner(ctx context.Context, task domain.Task) error {
	userID, err := currentUserID(ctx)
	if err != nil {
		return err
	}

	if task.AuthorUserID != userID {
		return fmt.Errorf("task with id=%d: %w", task.ID, core_errors.ErrNotFound)
	}

	return nil
}
