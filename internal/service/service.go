package service

import (
	"github.com/acool-kaz/bot-tasker/internal/repository"
)

type Service struct {
	User
	Task
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: newUserService(repo.User),
		Task: newTaskService(repo.Task),
	}
}
