package model

import (
	"gorm.io/gorm"
)

type IPBan struct {
	gorm.Model
	IP        string `gorm:"size:45;uniqueIndex" json:"ip"`
	Reason    string `gorm:"size:256" json:"reason"`
	BannedBy  uint   `json:"banned_by"`
	ExpiresAt *int64 `json:"expires_at"`
}

func GetAllIPBans() ([]IPBan, error) {
	var bans []IPBan
	err := DB.Find(&bans).Error
	return bans, err
}

func GetIPBansPaged(page, pageSize int) ([]IPBan, int64, error) {
	var bans []IPBan
	var total int64
	DB.Model(&IPBan{}).Count(&total)
	err := DB.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&bans).Error
	return bans, total, err
}

func (b *IPBan) Insert() error {
	return DB.Create(b).Error
}

func DeleteIPBan(id uint) error {
	return DB.Delete(&IPBan{}, id).Error
}
