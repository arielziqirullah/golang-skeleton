package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"golang/golang-skeleton/dto/pagination"
	"golang/golang-skeleton/dto/user"
	"golang/golang-skeleton/entity"
	"golang/golang-skeleton/helper"
	"golang/golang-skeleton/repository"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user user.UserUpdateRequestDTO) entity.User
	Profile(userID string) entity.User
	FindAll(user *entity.User, pagination *pagination.Pagination, search *user.SearchUser) (*[]entity.User, error)
	CountAll(search *user.SearchUser) (*int64, error)
	Import(file *multipart.FileHeader) error
	ExportGetData(users *entity.User, search *user.SearchUser) (*[]entity.User, error)
	Export(users *[]entity.User) (bytes.Buffer, error)
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

// service import file csv to database with bulk insert
func (service *userService) Import(file *multipart.FileHeader) error {
	csvFile, err := file.Open()
	if err != nil {
		return err
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	var users []entity.User
	for _, record := range records {
		user := entity.User{
			Name:      record[0],
			Email:     record[1],
			Password:  os.Getenv("DEFAULT_PASSWORD"),
			CreatedAt: time.Now(),
		}

		users = append(users, user)
	}

	service.userRepository.InsertUsers(users)

	return nil
}

func (service *userService) ExportGetData(users *entity.User, search *user.SearchUser) (*[]entity.User, error) {
	return service.userRepository.ExportGetData(users, search)
}

// service export file csv from database
func (service *userService) Export(users *[]entity.User) (bytes.Buffer, error) {
	var buf bytes.Buffer

	writer := csv.NewWriter(&buf)

	header := []string{"No", "Name", "Email"}
	writer.Write(header)

	usersSlice := *users

	startNumber := 1
	for _, user := range usersSlice {
		row := []string{strconv.Itoa(startNumber), user.Name, user.Email}
		writer.Write(row)
		startNumber++
	}

	writer.Flush()

	return buf, nil
}
