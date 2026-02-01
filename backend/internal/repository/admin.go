package repository

import (
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *AdminRepo { return &AdminRepo{db: db} }

func (r *AdminRepo) GetByUsername(username string) (*model.Admin, error) {
	var a model.Admin
	err := r.db.Where("username = ?", username).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepo) GetByID(id uint64) (*model.Admin, error) {
	var a model.Admin
	err := r.db.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepo) Create(a *model.Admin) error {
	return r.db.Create(a).Error
}

func (r *AdminRepo) ExistsByUsername(username string) (bool, error) {
	var c int64
	err := r.db.Model(&model.Admin{}).Where("username = ?", username).Count(&c).Error
	return c > 0, err
}
