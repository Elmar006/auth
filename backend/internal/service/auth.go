package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"auth/service/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  UserRepoService
	jwtSecret []byte
}

func NewAuthService(userRepo UserRepoService, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (string, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return "", errors.New("Email already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Failed to hashed password: %v", err)
		return "", errors.New("internal error")
	}

	user := model.Model{
		Name:     name,
		Email:    email,
		Password: string(hashPassword),
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Errorf("Failed to create user: %v", err)
		return "", errors.New("Failed to register")
	}

	token, err := s.generateJWT(id, email)
	if err != nil {
		log.Errorf("Failed to generate token: %v", err)
		return "", errors.New("Internal error")
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}
		log.Errorf("Failed to get user: %v", err)
		return "", errors.New("internal error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		log.Errorf("Failed to generate token: %v", err)
		return "", errors.New("internal error")
	}

	return token, nil
}

func (s *AuthService) generateJWT(userID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*model.Model, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := int64(claims["user_id"].(float64))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
