package exception

import (
	"golang/golang-skeleton/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("not route found")
		response := helper.BuildErrorResponse("not route found", http.StatusNotFound, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
	}
}
