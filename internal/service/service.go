package service

import (
	"homecloud--dbmanager-service/internal/interfaces"
)

type userService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
	return &userService{repo: repo}
}

// Реализация методов будет добавлена позже
