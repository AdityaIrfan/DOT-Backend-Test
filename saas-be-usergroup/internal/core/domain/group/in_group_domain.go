package group

import (
	"time"
)

type InGroup struct {
	ID            uint64
	UserGroupID   uint64
	UserAccountID uint64
	TimeAdded     time.Time
	TimeRemoved   time.Time
	GroupAdmin    bool
	Creator       bool
}

func NewInGroup() *InGroup {
	return &InGroup{}
}

func (i InGroup) IsGroupAdmin() bool {
	return i.GroupAdmin == true
}

func (i *InGroup) SetTimeAdded() {
	i.TimeAdded = time.Now()
}

func (i *InGroup) SetTimeRemoved() {
	i.TimeRemoved = time.Now()
}

func (i InGroup) ToHistory(historyType HistoryType) *InGroupHistory {
	return &InGroupHistory{
		InGroupID:     i.ID,
		UserAccountID: i.UserAccountID,
		Type:          historyType,
		HistoryTs:     time.Now(),
	}
}

func (i *InGroup) IsEmpty() bool {
	return i == nil
}

func (i InGroup) IsAdmin() bool {
	return i.GroupAdmin == true
}

func (i InGroup) IsCreator() bool {
	return i.Creator == true
}
