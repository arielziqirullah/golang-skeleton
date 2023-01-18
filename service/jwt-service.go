package service

import (
	"fmt"
	"golang/golang-skeleton/helper"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userID string, Email string) (string, *jwt.NumericDate)
	RefreshToken(UserID string, Email string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatalln("failed to load jwt secret key")
	}

	return secretKey
}

func (jwtservice *jwtService) GenerateToken(UserID string, Email string) (string, *jwt.NumericDate) {

	expToken := jwt.NewNumericDate(time.Now().Add(time.Minute * 15))

	claims := &jwtCustomClaim{
		UserID,
		jwt.RegisteredClaims{
			ExpiresAt: expToken,
			Subject:   Email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(jwtservice.secretKey))
	helper.LogIfError(err)

	return tokenSigned, expToken
}

func (jwtservice *jwtService) RefreshToken(UserID string, Email string) string {

	expToken := jwt.NewNumericDate(time.Now().Add(time.Hour * 24))

	claims := &jwtCustomClaim{
		UserID,
		jwt.RegisteredClaims{
			ExpiresAt: expToken,
			Subject:   Email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(jwtservice.secretKey))
	helper.LogIfError(err)

	return tokenSigned
}

func (jwtservice *jwtService) ValidateToken(token string) (*jwt.Token, error) {

	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(jwtservice.secretKey), nil
	})
}
