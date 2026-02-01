package repository

import (
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) GetByID(id uint64) (*model.User, error) {
	var u model.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByPhone(phone string) (*model.User, error) {
	var u model.User
	err := r.db.Where("phone = ?", phone).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(offset, limit int) ([]model.User, int64, error) {
	var list []model.User
	var total int64
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *UserRepo) Update(u *model.User) error {
	return r.db.Save(u).Error
}
