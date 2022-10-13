package group

import (
	"errors"
	"gorm.io/gorm"
)

func (u *UserGroup) Create(db *gorm.DB, userAccountID uint64) (*UserGroup, error) {
	if err := db.Transaction(func(tx *gorm.DB) error { // create user_group
		u.SetInsertTS()
		if err := tx.Save(&u).Error; err != nil {
			return err
		}

		// create in_group
		inGroup := u.ToInGroup(userAccountID)
		if err := tx.Save(&inGroup).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *UserGroup) Delete(db *gorm.DB) error {
	if err := db.Delete(&u).Error; err != nil {
		return err
	}

	return nil
}

func (i *InGroup) GetOneByUserGroupIDAndUserAccountID(db *gorm.DB, userGroupID uint64, UserAccountID uint64) (*InGroup, error) {
	if err := db.Where("user_group_id = ?", userGroupID).Where("user_account_id = ?", UserAccountID).First(&i).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return i, nil
}

func (i *InGroup) AddMember(db *gorm.DB, adminID uint64) (*InGroup, error) {
	if err := db.Transaction(func(tx *gorm.DB) error {
		// create in_groyp
		i.SetTimeAdded()
		if err := tx.Save(&i).Error; err != nil {
			return err
		}

		// create in_group_history
		history := i.ToHistory(HistoryAdded).SetAdminID(adminID)
		if err := tx.Save(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *InGroup) RemoveMember(db *gorm.DB, adminID uint64) (*InGroup, error) {
	if err := db.Transaction(func(tx *gorm.DB) error {
		i.SetTimeRemoved()

		// update for time_removed
		if err := tx.Save(&i).Error; err != nil {
			return err
		}

		// create history
		history := i.ToHistory(HistoryRemoved).SetAdminID(adminID)
		if err := tx.Save(&history).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return i, nil
}

func (u *UserGroup) GetOneByID(db *gorm.DB, id uint64) (*UserGroup, error) {
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (u *UserGroupType) GetOneByID(db *gorm.DB, id uint64) (*UserGroupType, error) {
	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (u *InGroup) GetTotalMember(db *gorm.DB, userGroupID uint64) (int, error) {
	var totalMember *int64

	if err := db.Where("user_group_id = ?", userGroupID).Where("time_removed IS NULL").Find(&u).Count(totalMember).Error; err != nil {
		return 0, err
	}

	return int(*totalMember), nil
}

// params bool on first parameter means isErrorInternal
func (i *InGroup) IsAvailableToAddOneMore(db *gorm.DB) (bool, error) {
	userGroup, err := NewUserGroup().GetOneByID(db, i.UserGroupID)
	if err != nil {
		return true, err
	}

	if userGroup.IsEmpty() {
		return false, errors.New("user group not found")
	}

	userGroupType, err := NewUserGroupType().GetOneByID(db, userGroup.UserGroupTypeID)
	if err != nil {
		return true, err
	}

	if userGroupType.IsEmpty() {
		return false, errors.New("user group type not found")
	}

	totalMember, err := i.GetTotalMember(db, i.UserGroupID)
	if err != nil {
		return true, err
	}

	if !userGroupType.IsAvailableToAddOneMore(totalMember) {
		return false, errors.New("can not add member because the slot is empty")
	}

	return false, nil
}
