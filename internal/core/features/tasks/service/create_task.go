package tasks_service

import (
	"context"
	"fmt"

	"github.com/stfprod17/todoapp/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get current user: %w", err)
	}

	task.AuthorUserID = userID

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	task, err = s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil

}
