package service

import (
	"fmt"
	"golang/golang-skeleton/dto/authentication"
	"golang/golang-skeleton/entity"
	"golang/golang-skeleton/helper"
	"golang/golang-skeleton/repository"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user authentication.RegisterRequestDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRep repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRep,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	if value, ok := res.(entity.User); ok {
		comparedPassword := comparedPassword(value.Password, []byte(password))
		if value.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authService) CreateUser(user authentication.RegisterRequestDTO) entity.User {

	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	helper.LogIfError(fmt.Errorf("failed map %v", err))

	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authService) FindByEmail(email string) entity.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparedPassword(hashPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		helper.LogIfError(err)
		return false
	}
	return true
}
