package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/api/session"
	"apps/center/db"
	"apps/center/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Login 登录
func Login(c *gin.Context) {

	loginForm := param.LoginForm{}

	err := c.ShouldBindJSON(&loginForm)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	userInfo := db.UserInfo{}

	tx := global.DB.Where("account = ?", loginForm.Account).Find(&userInfo)

	if tx.RowsAffected == 0 || tx.Error != nil {
		c.JSON(http.StatusOK, response.Resp[any]{
			Code: 403,
			Msg:  "账号或密码错误",
			Data: nil,
		})
		return
	}

	account := userInfo.Account
	password := userInfo.Password

	if account != loginForm.Account {
		c.JSON(http.StatusOK, response.Resp[any]{
			Code: 403,
			Msg:  "账号或密码错误",
			Data: nil,
		})
		return
	}

	if password != loginForm.PassWord {
		c.JSON(http.StatusOK, response.Resp[any]{
			Code: 403,
			Msg:  "账号或密码错误",
			Data: nil,
		})
		return
	}

	//新建session
	newSession := session.NewSession()

	newSession.Set("account", account)
	newSession.Set("password", password)
	newSession.Set("Address", c.Request.RemoteAddr)
	newSession.Set("time", time.Now())

	c.JSON(http.StatusOK, response.Success(newSession.Id))

}

func CheckLogin(c *gin.Context) {

	c.JSON(http.StatusOK, response.Success[any](nil))

}
