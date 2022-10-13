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

func (u *User) GetOneByEmail(db *gorm.DB, email string) (*User, error) {
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (u *User) GetOneByUUID(db *gorm.DB, uuid string) (*User, error) {
	if err := db.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (u *User) GetOneByUsername(db *gorm.DB, username string) (*User, error) {
	if err := db.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (u *User) IsUsernameAvailable(db *gorm.DB, username string) (bool, error) {
	u, err := u.GetOneByUsername(db, username)
	if err != nil {
		return false, err
	}

	if u.IsEmpty() {
		return true, nil
	}

	return false, nil
}

func (u *User) IsEmailAvailable(db *gorm.DB, email string) (bool, error) {
	u, err := u.GetOneByEmail(db, email)
	if err != nil {
		return false, err
	}

	if u.IsEmpty() || !u.IsVerified() {
		return true, nil
	}

	return false, nil
}
