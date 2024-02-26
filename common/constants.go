// Package common
/**
  @author: zhenxiangjin
  @update_date: 2023/12/6
  @note:
*/
package common

import (
	"github.com/google/uuid"
	"time"
)

var Version = "v0.0.0" // this hard coding will be replaced automatically when building, no need to manually change
var SystemName = "CoolRun"
var ServerAddress = "http://localhost:3000"
var Footer = ""
var HomePageLink = ""

var SessionSecret = uuid.New().String()

var WeChatAppID = "xxx"
var WeChatAppSecret = "xxxx"
var WeChatServerAddress = "xxx"
var WeChatServerToken = ""
var WeChatAccountQRCodeImageURL = ""

var WeChatAuthEnabled = true
var RegisterEnabled = true

var RateLimitKeyExpirationDuration = 20 * time.Minute

// All duration's unit is seconds
// Shouldn't larger than RateLimitKeyExpirationDuration
var (
	GlobalApiRateLimitNum            = 60
	GlobalApiRateLimitDuration int64 = 3 * 60

	GlobalWebRateLimitNum            = 60
	GlobalWebRateLimitDuration int64 = 3 * 60

	UploadRateLimitNum            = 10
	UploadRateLimitDuration int64 = 60

	DownloadRateLimitNum            = 10
	DownloadRateLimitDuration int64 = 60

	CriticalRateLimitNum            = 20
	CriticalRateLimitDuration int64 = 20 * 60
)

const (
	RoleGuestUser  = 0
	RoleCommonUser = 1
	RoleAdminUser  = 10
	RoleRootUser   = 100
)

const (
	UserStatusEnabled  = 1 // don't use 0, 0 is the default value!
	UserStatusDisabled = 2 // also don't use 0
)

const (
	RoomCapacity    = 5
	VIPRoomCapacity = 10
)
