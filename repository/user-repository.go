package repository

import (
	"fmt"
	"golang/golang-skeleton/dto/pagination"
	"golang/golang-skeleton/dto/user"
	"golang/golang-skeleton/entity"
	"golang/golang-skeleton/helper"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) entity.User
	InsertUsers(users []entity.User) []entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	ProfileUser(userID string) entity.User
	FindAll(user *entity.User, pagination *pagination.Pagination, search *user.SearchUser) (*[]entity.User, error)
	CountAll(search *user.SearchUser) (*int64, error)
	ExportGetData(user *entity.User, search *user.SearchUser) (*[]entity.User, error)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user
}

func (db *userConnection) InsertUsers(users []entity.User) []entity.User {
	for _, user := range users {
		user.Password = hashAndSalt([]byte(user.Password))
		db.connection.Save(&user)
	}
	return users
}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tmpUser entity.User
		db.connection.Find(&tmpUser, user.ID)
		user.Password = tmpUser.Password
	}

	db.connection.Save(&user)
	return user
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}

	return nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) entity.User {
	var user entity.User
	db.connection.Where("email = ?", email).Take(&user)
	return user
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User
	db.connection.Find(&user, userID)
	return user
}

func (db *userConnection) FindAll(user *entity.User, pagination *pagination.Pagination, search *user.SearchUser) (*[]entity.User, error) {
	var users []entity.User

	offset := (pagination.Page - 1) * pagination.Limit
	queryBuider := db.connection.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	if search.SearchName != "" {
		queryBuider.Where("name LIKE ?", "%"+search.SearchName+"%")
	}
	result := queryBuider.Model(&entity.User{}).Where(user).Find(&users)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return &users, nil
}

func (db *userConnection) ExportGetData(user *entity.User, search *user.SearchUser) (*[]entity.User, error) {
	var users []entity.User

	queryBuider := db.connection
	if search.SearchName != "" {
		queryBuider.Where("name LIKE ?", "%"+search.SearchName+"%")
	}
	result := queryBuider.Model(&entity.User{}).Where(user).Find(&users)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return &users, nil
}

func (db *userConnection) CountAll(search *user.SearchUser) (*int64, error) {
	var count int64

	queryBulder := db.connection.Model(&entity.User{})
	if search.SearchName != "" {
		queryBulder.Where("name LIKE ?", "%"+search.SearchName+"%")
	}
	result := queryBulder.Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}
	return &count, nil

}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	helper.LogIfError(fmt.Errorf("failed to hash password, Error : %s", err))

	return string(hash)
}
