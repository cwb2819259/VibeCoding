package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/repository"
	"gorm.io/gorm"
)

type AuthUserService struct {
	userRepo *repository.UserRepo
	secret   string
	expire   time.Duration
}

func NewAuthUserService(userRepo *repository.UserRepo, secret string, expireHours int) *AuthUserService {
	return &AuthUserService{
		userRepo: userRepo,
		secret:   secret,
		expire:   time.Duration(expireHours) * time.Hour,
	}
}

const MockCode = "123456"

func (s *AuthUserService) Login(phone, code string) (token string, user *model.User, err error) {
	if code != MockCode {
		return "", nil, errors.New("invalid code")
	}
	u, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u = &model.User{Phone: phone, Nickname: "用户" + phone[len(phone)-4:]}
			if err = s.userRepo.Create(u); err != nil {
				return "", nil, err
			}
		} else {
			return "", nil, err
		}
	}
	claims := &middleware.UserClaims{
		UserID: u.ID,
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
	return token, u, nil
}
