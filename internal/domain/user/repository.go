package user

import "github.com/felipeversiane/go-boiterplate/internal/infra/database"

type userRepository struct {
	db database.DatabaseInterface
}

type UserRepositoryInterface interface {
}

func NewUserRepository(db database.DatabaseInterface) UserRepositoryInterface {
	return &userRepository{db}
}
