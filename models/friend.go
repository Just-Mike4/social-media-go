package models

import "gorm.io/gorm"

type FriendRequest struct {
	gorm.Model
	FromUserID uint
	ToUserID   uint
	Status     string `json:"status"` // "pending", "accepted", "rejected"
}
