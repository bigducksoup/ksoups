package api

import (
	"config-manager/center/apiserver/params"
	"config-manager/center/apiserver/response"
	"config-manager/center/apiserver/session"
	"config-manager/center/config"
	"config-manager/common/utils"
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
	"time"
)

func Login(c *gin.Context) {

	loginForm := params.LoginForm{}

	err := c.ShouldBindJSON(&loginForm)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	account := config.Conf.Account
	password := config.Conf.Password

	if account != loginForm.Account {
		c.JSON(http.StatusOK, response.Resp[any]{
			Code: 403,
			Msg:  "账号或密码错误",
			Data: nil,
		})
		return
	}

	originPassMd5 := utils.Md5([]byte(password))

	if originPassMd5 != loginForm.PassWord {
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
