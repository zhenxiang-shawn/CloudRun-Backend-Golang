// Package model
/**
  @author: zhenxiangjin
  @update_date: 2023/12/13
  @note:
*/
package model

import (
	"coolrun-backend-go/common"
	"errors"
)

type User struct {
	Model
	Username         string `json:"username" gorm:"unique;index" validate:"max=12"`
	Password         string `json:"password" gorm:"not null;" validate:"min=8,max=20"`
	DisplayName      string `json:"display_name" gorm:"index" validate:"max=20"`
	Role             int    `json:"role" gorm:"type:int;default:1"`   // admin, common
	Status           int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Token            string `json:"token" gorm:"index"`
	Email            string `json:"email" gorm:"index" validate:"max=50"`
	WeChatId         string `json:"wechat_id" gorm:"column:wechat_id;index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"` // this field is only for Email verification, don't save it to database!
}

// CleanUser
// This struct is for removing sensitive information before
// sending back information to client
type CleanUser struct {
	ID               int    `json:"id"`
	Username         string `json:"username" `
	DisplayName      string `json:"display_name" gorm:"index" validate:"max=20"`
	Role             int    `json:"role" gorm:"type:int;default:1"` // admin, common
	Email            string `json:"email" gorm:"index" validate:"max=50"`
	WeChatId         string `json:"wechat_id" gorm:"column:wechat_id;index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"` // this fiel
}

func GetMaxUserID() uint {
	var user User
	// TODO(zhenxiang) resolve this: Will return error when the dataset is empty
	DB.Last(&user)
	return user.Model.ID
}

func IsEmailAlreadyTaken(email string) bool {
	return DB.Where("email = ?", email).Find(&User{}).RowsAffected == 1
}

func IsWeChatIdAlreadyTaken(wechatId string) bool {
	return DB.Where("wechat_id = ?", wechatId).Find(&User{}).RowsAffected == 1
}

func IsUsernameAlreadyTaken(username string) bool {
	return DB.Where("username = ?", username).Find(&User{}).RowsAffected == 1
}

func (user *User) Insert() error {
	var err error
	if user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	err = DB.Create(user).Error
	return err
}

func (user *User) Update() error {
	var err error
	err = DB.Model(user).Updates(user).Error
	return err
}

func (user *User) FillUserById() error {
	if user.ID == 0 {
		return errors.New("id 为空！")
	}
	DB.Where(User{Model: Model{ID: user.ID}}).First(user)
	return nil
}

func (user *User) FillUserByWeChatId() error {
	if user.WeChatId == "" {
		return errors.New("WeChat id 为空！")
	}
	DB.Where(User{WeChatId: user.WeChatId}).First(user)
	return nil
}

func (user *User) FillUserByUsername() error {
	if user.Username == "" {
		return errors.New("username 为空！")
	}
	DB.Where(User{Username: user.Username}).First(user)
	return nil
}

func (user *User) CleanSensitiveInfo() error {
	user.Email = ""
	user.Password = ""
	user.Token = ""
	user.WeChatId = ""
	return nil
}

/*
	cleanUser := model.User{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		Status:      user.Status,
	}
*/
