package services

import (
	"errors"
	"meetingassist/config"
	"meetingassist/models"
	"meetingassist/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	db := config.GetDB()

	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	db := config.GetDB()

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiration)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}
