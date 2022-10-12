package user

import (
	"errors"
	"gorm.io/gorm"
)

func (u *User) Create(db *gorm.DB) (*User, error) {
	if err := db.Save(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Update(db *gorm.DB) (*User, error) {
	if err := db.Save(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) GetOneByID(db *gorm.DB, id uint64) (*User, error) {
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (u *User) GetOneByEmail(db *gorm.DB) (*User, error) {
	if err := db.Where("email = ?", u.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (u *User) GetOneUUID(db *gorm.DB) (*User, error) {
	if err := db.Where("uuid = ?", u.UUID).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}
