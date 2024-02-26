// Package controller
/**
  @author: zhenxiangjin
  @update_date: 2023/12/15
  @note:
*/
package controller

import (
	"coolrun-backend-go/common"
	"coolrun-backend-go/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

// 封装code2session接口返回数据
type wechatLoginResponse struct {
	ErrMsg     string `json:"errmsg"`
	ErrCode    int    `json:"errcode"`
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
}

func getWeChatIdByCode(code string) (string, error) {
	if code == "" {
		return "", errors.New("无效的参数")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", common.WeChatServerAddress, common.WeChatAppID, common.WeChatAppSecret, code), nil)

	if err != nil {
		common.SysError(err.Error())
		return "", err
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	httpResponse, err := client.Do(req)
	body, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		fmt.Printf("something wrong1\n")
		return "", err
	}
	defer httpResponse.Body.Close()
	var res wechatLoginResponse
	err = json.Unmarshal(body, &res)

	// check the error code and return result based on error code
	switch res.ErrCode {
	case 0:
		return res.OpenID, nil
	case 40029:
		return "", errors.New("验证码错误或已过期")
	case 15011:
		return "", errors.New("api minute-quota reach limit must slower retry next minute")
	case 40226:
		return "", errors.New("code blocked")
	case -1:
		return "", errors.New("system error")
	default:
		return "", errors.New("unknown error")
	}
}

// WechatAuth 微信登录验证
// @Summary 微信小程序登录/注册
// @Param code query string true "wx.Login request code"
// @Success 200 {object} model.CleanUser "成功"
// @Router /login/WeChat [get]
func WechatAuth(c *gin.Context) {
	code := c.Query("code")
	//fmt.Printf("code: %s\n", code)
	wechatId, err := getWeChatIdByCode(code)
	//fmt.Printf("wechatId: %s,\n err: %v\n", wechatId, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// check WeChat ID is taken or not
	user := model.User{
		WeChatId: wechatId,
	}
	if model.IsWeChatIdAlreadyTaken(wechatId) {
		err := user.FillUserByWeChatId()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			fmt.Printf("x\n")
			common.SysError(err.Error())
			return
		}
		cleanUser := model.User{
			Model: model.Model{
				ID: user.ID,
			},
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Role:        user.Role,
			Status:      user.Status,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User has been signed in.",
			"data":    cleanUser,
		})

	} else {
		if common.RegisterEnabled {
			userId := strconv.Itoa(int(model.GetMaxUserID() + 1))
			user.Username = "wechat_" + userId
			user.DisplayName = "微信用户" + userId
			user.Role = common.RoleCommonUser
			user.Status = common.UserStatusEnabled

			if err := user.Insert(); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": err.Error(),
				})
				common.SysError(err.Error())
				return
			}
			if err != nil {
				cleanUser := model.User{
					Model: model.Model{
						ID: user.ID,
					},
					Username:    user.Username,
					DisplayName: user.DisplayName,
					Role:        user.Role,
					Status:      user.Status,
				}
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "User has been signed up.",
					"data":    cleanUser,
				})
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "管理员关闭了新用户注册",
			})
			return
		}
	}

}
