package service

import (
	backend "cmd/main.go"
	"cmd/main.go/pkg/repository"
	"crypto/sha1"
	"fmt"
	"time"



	"github.com/dgrijalva/jwt-go"
)
const (
	salt = "fjsdlkfjdslkfnvmxcfjkdslf"
	tokenTTL   = time.Hour * 12 //время сколько работает токен
	signingKey = "123dfjdsk123dfsk$L#FJK"
)
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

type AuthService struct {
	repo repository.Authorization
}

func (s *AuthService) CreateUser(user backend.User) (int, error){
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error){
	user,err := s.repo.GetUser(username, generatePasswordHash(password))
	if err!= nil{
		fmt.Println(err.Error())
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}
func (s *AuthService) ParseToken(accesToken string) (int, error){
	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0,  fmt.Errorf("token claims are not of type *tokenClaims ")
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
