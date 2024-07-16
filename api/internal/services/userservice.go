package services

import (
	"em-test/internal/adapters"
	"em-test/internal/domain"
	"em-test/internal/lib/dto"
	"em-test/internal/lib/filters"
	"log/slog"
)

var _ adapters.UsersService = (*UsersService)(nil)

type UserRepository interface {
	Add(dto dto.SaveUserDto) (*domain.User, error)
	Read(id string) (*domain.User, error)
	ReadMany(filters *filters.UsersFilters) ([]*domain.User, int64, error)
}

type UserFinder interface {
	GetInfo(serie int, number int) (*dto.UserInfoDto, error)
}

type UsersService struct {
	repository UserRepository
	userFinder UserFinder
}

func NewUserService(userRepository UserRepository, passportApiRepository UserFinder) *UsersService {
	return &UsersService{
		repository: userRepository,
		userFinder: passportApiRepository,
	}
}

func (u *UsersService) AddUser(addUserDto *dto.AddUserDto) (*domain.User, error) {
	const fn = "UsersService.AddUser"
	logger := slog.With(slog.String("fn", fn))

	info, err := u.userFinder.GetInfo(addUserDto.PassportSerie, addUserDto.PassportNumber)
	if err != nil {
		logger.Error("user not found", slog.String("err", err.Error()))
		return nil, err
	}
	logger.Debug("user found", slog.Any("user", info))

	saveUserDto := dto.SaveUserDto{
		AddUserDto:  addUserDto,
		UserInfoDto: info,
	}

	user, err := u.repository.Add(saveUserDto)
	if err != nil {
		logger.Error("error with saving user in repository", slog.Any("save user dto", saveUserDto), slog.String("err", err.Error()))
		return nil, err
	}
	logger.Debug("user saved", slog.Any("user", user))

	return user, nil
}

func (u *UsersService) GetUsers(filters *filters.UsersFilters) (users []*domain.User, total int64, err error) {
	const fn = "UsersService.GetUsers"
	logger := slog.With(slog.String("fn", fn))

	logger.Debug("get users", slog.Any("filters", filters))

	users, total, err = u.repository.ReadMany(filters)
	if err != nil {
		logger.Error("error with getting users from repository", slog.Any("filters", filters), slog.String("err", err.Error()))
		return nil, 0, err
	}

	logger.Debug("got users", slog.Any("users", users))

	return users, total, nil
}
