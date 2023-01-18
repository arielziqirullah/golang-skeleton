package controller

import (
	"fmt"
	"golang/golang-skeleton/dto/authentication"
	"golang/golang-skeleton/entity"
	"golang/golang-skeleton/helper"
	"golang/golang-skeleton/service"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthController interface {
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authservice service.AuthService, jwtservice service.JWTService) AuthController {
	return &authController{
		authService: authservice,
		jwtService:  jwtservice,
	}
}

func (auth *authController) Login(ctx *gin.Context) {
	var loginDTO authentication.LoginRequestDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		log.Println(errDTO.Error())
		response := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authResult := auth.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)

	var loginResponseDTO authentication.LoginResponseDTO
	if value, ok := authResult.(entity.User); ok {
		// Generate Token JWT
		generateToken, expToken := auth.jwtService.GenerateToken(strconv.FormatUint(value.ID, 10), value.Email)
		loginResponseDTO.AccessToken = generateToken
		loginResponseDTO.ExpiresIn = int(expToken.Unix())
		loginResponseDTO.TokenType = os.Getenv("TOKEN_TYPE")
		response := helper.BuildResponse(true, "OK", http.StatusOK, loginResponseDTO)

		// Generate Refresh Token
		refreshToken := auth.jwtService.RefreshToken(strconv.FormatUint(value.ID, 10), value.Email)
		refName := "refresh_token"
		refValue := refreshToken
		refPath := "/"
		refDomain := os.Getenv("APP_RUN")
		refSecure := false
		refHttp := true

		refMaxAge, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
		helper.LogIfError(err)

		ctx.SetCookie(refName, refValue, refMaxAge, refPath, refDomain, refSecure, refHttp)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("please check again your credential", http.StatusUnauthorized, helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (auth *authController) RefreshToken(ctx *gin.Context) {

	var refreshToken authentication.RefreshTokenRequestDTO
	errRef := ctx.ShouldBind(&refreshToken)
	if errRef != nil {
		helper.LogIfError(errRef)
		response := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	validateToken, err := auth.jwtService.ValidateToken(refreshToken.TokenRefresh)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("failed to validate token", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	claims := validateToken.Claims.(jwt.MapClaims)
	email := fmt.Sprintf("%v", claims["sub"])

	checkEmail := auth.authService.FindByEmail(email)
	if checkEmail.Email == "" {
		log.Println("failed to find email")
		response := helper.BuildErrorResponse("failed to find email", http.StatusNotFound, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	id, errParse := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if errParse == nil {
		helper.LogIfError(err)

		var loginResponseDTO authentication.LoginResponseDTO
		generateToken, expToken := auth.jwtService.GenerateToken(strconv.FormatUint(id, 10), email)
		loginResponseDTO.AccessToken = generateToken
		loginResponseDTO.ExpiresIn = int(expToken.Unix())
		loginResponseDTO.TokenType = os.Getenv("TOKEN_TYPE")
		response := helper.BuildResponse(true, "OK", http.StatusOK, loginResponseDTO)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("please check again your credential", http.StatusUnauthorized, helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (auth *authController) Register(ctx *gin.Context) {
	var registerDTO authentication.RegisterRequestDTO

	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		log.Println(errDTO.Error())
		response := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !auth.authService.IsDuplicateEmail(registerDTO.Email) {
		log.Println("duplicate email")
		response := helper.BuildErrorResponse("failed to process request", http.StatusConflict, helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser := auth.authService.CreateUser(registerDTO)
		response := helper.BuildResponse(true, "OK", http.StatusCreated, createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
