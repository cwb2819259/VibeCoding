package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type AuthAdminService struct {
	adminRepo *repository.AdminRepo
	secret    string
	expire    time.Duration
}

func NewAuthAdminService(adminRepo *repository.AdminRepo, secret string, expireHours int) *AuthAdminService {
	return &AuthAdminService{
		adminRepo: adminRepo,
		secret:    secret,
		expire:    time.Duration(expireHours) * time.Hour,
	}
}

func (s *AuthAdminService) Login(username, password string) (token string, admin *model.Admin, err error) {
	a, err := s.adminRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid username or password")
		}
		return "", nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid username or password")
	}
	claims := &middleware.AdminClaims{
		AdminID: a.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(s.secret))
	if err != nil {
		return "", nil, err
	}
	return token, a, nil
}
