package routes

import (
	"golang/golang-skeleton/config"
	"golang/golang-skeleton/exception"
	"golang/golang-skeleton/middleware"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Run() {

	defer config.CloseDatabaseConnection(db)
	config.AppLoadEnv()
	config.AppDebug()

	routes := gin.Default()

	routes.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":         "OK",
			"message":        "Golang Skeleton",
			"golang_version": runtime.Version(),
		})
	})

	// Register route
	routesV1 := routes.Group("v1")
	{
		authRoutes := routesV1.Group("/api/auth")
		{
			authRoutes.POST("/login", authController.Login)
			authRoutes.POST("/refresh", authController.RefreshToken)
			authRoutes.POST("/register", authController.Register)
		}

		userRoutes := routesV1.Group("/api/user", middleware.AuthorizeJWT(jwtService))
		{
			userRoutes.GET("/profile", userController.Profile)
			userRoutes.GET("/find-all", userController.FindAll)
			userRoutes.PUT("/profile", userController.Update)
			userRoutes.POST("/import", userController.Import)
		}
	}

	routes.NoRoute(exception.NotFoundRoute())
	if os.Getenv("RUNNING_PORT") == "" {
		log.Println("listening and serving on default port :8080")
		routes.Run()
	} else {
		log.Println("listening and serving on port " + os.Getenv("RUNNING_PORT"))
		routes.Run(os.Getenv("RUNNING_PORT"))
	}
}
