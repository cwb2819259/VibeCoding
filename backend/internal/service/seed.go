package service

import (
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminIfNotExists(adminRepo *repository.AdminRepo) error {
	exists, err := adminRepo.ExistsByUsername("admin")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a := &model.Admin{
		Username:     "admin",
		PasswordHash: string(hash),
		Nickname:     "Administrator",
	}
	return adminRepo.Create(a)
}
