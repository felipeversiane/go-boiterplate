package user

type userUsecase struct {
	repository UserRepositoryInterface
}

type UserUsecaseInterface interface {
}

func NewUserUsecase(repository UserRepositoryInterface) UserUsecaseInterface {
	return &userUsecase{repository}
}
