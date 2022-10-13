package group

import (
	"time"
)

type HistoryType string

const (
	HistoryAdded   HistoryType = "added"
	HistoryRemoved HistoryType = "removed"
)

type InGroupHistory struct {
	InGroupID     uint64
	AdminID       uint64
	UserAccountID uint64
	Type          HistoryType
	HistoryTs     time.Time
}

func (i *InGroupHistory) SetAdminID(adminID uint64) *InGroupHistory {
	i.AdminID = adminID
	return i
}
