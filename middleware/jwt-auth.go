package middleware

import (
	"golang/golang-skeleton/helper"
	"golang/golang-skeleton/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		authHeader = strings.Split(authHeader, "Bearer ")[1]
		if authHeader == "" {
			response := helper.BuildErrorResponse("failed to process request", http.StatusBadRequest, nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		}

		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[user_id]: ", claims["user_id"])
		} else {
			log.Println(err)
			response := helper.BuildErrorResponse("token is not valid", http.StatusUnauthorized, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
