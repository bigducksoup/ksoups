package handler

import (
	"config-manager/center/api/param"
	"config-manager/center/api/response"
	"config-manager/center/api/session"
	"config-manager/center/db"
	"config-manager/center/global"
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
	"time"
)

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
		c.JSON(http.StatusOK, response.Resp[types.Nil]{
			Code: 403,
			Msg:  "账号或密码错误",
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
