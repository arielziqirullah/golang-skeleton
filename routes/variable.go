package routes

import (
	"golang/golang-skeleton/config"
	"golang/golang-skeleton/controller"
	"golang/golang-skeleton/repository"
	"golang/golang-skeleton/service"

	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetUpDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
)
