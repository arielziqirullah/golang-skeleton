package service

import (
	"fmt"
	"golang/golang-skeleton/dto/pagination"
	"golang/golang-skeleton/dto/user"
	"golang/golang-skeleton/entity"
	"golang/golang-skeleton/helper"
	"golang/golang-skeleton/repository"

	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user user.UserUpdateRequestDTO) entity.User
	Profile(userID string) entity.User
	FindAll(user *entity.User, pagination *pagination.Pagination, search *user.SearchUser) (*[]entity.User, error)
	CountAll(search *user.SearchUser) (*int64, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user user.UserUpdateRequestDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	helper.LogIfError(fmt.Errorf("failed to map %s", err))

	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}

func (service *userService) FindAll(user *entity.User, pagination *pagination.Pagination, search *user.SearchUser) (*[]entity.User, error) {
	return service.userRepository.FindAll(user, pagination, search)
}

func (service *userService) CountAll(search *user.SearchUser) (*int64, error) {
	return service.userRepository.CountAll(search)
}
