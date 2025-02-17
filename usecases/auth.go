package usecases

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lucastomic/snoopycoon_back/database"
	"github.com/lucastomic/snoopycoon_back/domain"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user domain.User) (string, error) {
	u, _ := findUserByEmail(user.Email)
	if u.ID != 0 {
		return "", errors.New("user already exists")
	}
	encodedPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = encodedPassword
	err = database.DB.Create(&user).Error
	if err != nil {
		return "", err
	}
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func SignIn(email string, password string) (string, error) {
	user, err := findUserByEmail(email)
	if err != nil {
		return "", err
	}
	passwordIsCorrect := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passwordIsCorrect != nil {
		return "", errors.New("wrong password")
	}
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateJWT(user domain.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	idString := strconv.Itoa(int(user.ID))
	claims := &jwt.StandardClaims{
		Subject:   idString,
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func findUserByEmail(email string) (domain.User, error) {
	var user domain.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
