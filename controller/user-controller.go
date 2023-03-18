package controller

import (
	"fmt"
	"golang/golang-skeleton/dto/pagination"
	"golang/golang-skeleton/dto/user"
	"golang/golang-skeleton/entity"
	"golang/golang-skeleton/helper"
	"golang/golang-skeleton/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController interface {
	Update(ctx *gin.Context)
	Profile(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Import(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userservice service.UserService, jwtservice service.JWTService) UserController {
	return &userController{
		userService: userservice,
		jwtService:  jwtservice,
	}
}

func (userController *userController) Update(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	authHeader = strings.Split(authHeader, "Bearer ")[1]
	token, errToken := userController.jwtService.ValidateToken(authHeader)
	helper.LogIfError(errToken)

	var userUpdateDTO user.UserUpdateRequestDTO
	errDTO := ctx.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		log.Println(errDTO.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	helper.LogIfError(err)

	userUpdateDTO.ID = id
	userUpdateDTO.UpdatedAt = time.Now()
	userToUpdate := userController.userService.Update(userUpdateDTO)
	response := helper.BuildResponse(true, "OK", http.StatusOK, userToUpdate)
	ctx.JSON(http.StatusOK, response)

}

func (userController *userController) Profile(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	authHeader = strings.Split(authHeader, "Bearer ")[1]
	token, err := userController.jwtService.ValidateToken(authHeader)
	helper.LogIfError(err)

	claims := token.Claims.(jwt.MapClaims)
	user := userController.userService.Profile(fmt.Sprintf("%v", claims["user_id"]))
	response := helper.BuildResponse(true, "OK", http.StatusOK, user)
	ctx.JSON(http.StatusOK, response)
}

func (userController *userController) FindAll(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	authHeader = strings.Split(authHeader, "Bearer ")[1]
	_, err := userController.jwtService.ValidateToken(authHeader)
	helper.LogIfError(err)

	var Pagination pagination.PaginationRequest

	errPag := ctx.ShouldBind(&Pagination)
	if errPag != nil {
		log.Println(errPag.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	getPaginate := helper.GeneratePaginationFromRequest(&Pagination)

	var search user.SearchUser
	errSearch := ctx.ShouldBind(&search)
	if errSearch != nil {
		log.Println(errSearch.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var users entity.User
	userList, err := userController.userService.FindAll(&users, &getPaginate, &search)
	if err != nil {
		log.Println(err.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	getTotal, err := userController.userService.CountAll(&search)
	if err != nil {
		log.Println(err.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	metadataBuild := helper.BuildMetadataResponse(*getTotal, getPaginate.Limit, getPaginate.Page)

	var dataResponse pagination.DataResponse
	dataResponse.Data = userList
	dataResponse.Metadata = metadataBuild

	response := helper.BuildResponse(true, "OK", http.StatusOK, dataResponse)
	ctx.JSON(http.StatusOK, response)
}

// api import csv file to database
func (userController *userController) Import(ctx *gin.Context) {
	log.Println("[start controller | userController.Import]")
	authHeader := ctx.GetHeader("Authorization")
	authHeader = strings.Split(authHeader, "Bearer ")[1]
	_, err := userController.jwtService.ValidateToken(authHeader)
	helper.LogIfError(err)

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userController.userService.Import(file)
	response := helper.BuildResponse(true, "OK", http.StatusOK, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

// export csv file to api
func (userController *userController) Export(ctx *gin.Context) {
	log.Println("[start controller | userController.Export]")
	authHeader := ctx.GetHeader("Authorization")
	authHeader = strings.Split(authHeader, "Bearer ")[1]
	_, err := userController.jwtService.ValidateToken(authHeader)
	helper.LogIfError(err)

	var search user.SearchUser
	errSearch := ctx.ShouldBind(&search)
	if errSearch != nil {
		log.Println(errSearch.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var users entity.User
	userList, err := userController.userService.ExportGetData(&users, &search)
	if err != nil {
		log.Println(err.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	buf, err := userController.userService.Export(userList)
	if err != nil {
		log.Println(err.Error())
		res := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Set the response headers
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=users.csv")
	ctx.Header("Content-Type", "text/csv")

	// Write the CSV data to the response
	ctx.Data(http.StatusOK, "text/csv", buf.Bytes())
}
