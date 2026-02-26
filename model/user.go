package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LinuxDOID   int    `gorm:"uniqueIndex;column:linux_do_id" json:"linux_do_id"`
	Username    string `gorm:"size:64;uniqueIndex" json:"username"`
	DisplayName string `gorm:"size:128" json:"display_name"`
	AvatarURL   string `gorm:"size:512" json:"avatar_url"`
	TrustLevel  int    `json:"trust_level"`
	Role        int    `gorm:"default:1" json:"role"`
	Status      int    `gorm:"default:1" json:"status"`
	QuotaTotal  int64  `gorm:"default:1000" json:"quota_total"`
	QuotaUsed   int64  `gorm:"default:0" json:"quota_used"`
	TokenLimit  int    `gorm:"default:5" json:"token_limit"`
	LastLoginAt *int64 `json:"last_login_at"`
	LastLoginIP string `gorm:"size:45" json:"last_login_ip"`
}

func GetUserByLinuxDOID(id int) (*User, error) {
	var user User
	err := DB.Where("linux_do_id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64
	DB.Model(&User{}).Count(&total)
	err := DB.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (u *User) Insert() error {
	return DB.Create(u).Error
}

func (u *User) Update() error {
	return DB.Save(u).Error
}

func GetUserCount() int64 {
	var count int64
	DB.Model(&User{}).Count(&count)
	return count
}
