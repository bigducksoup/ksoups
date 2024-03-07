package handler

import (
	"apps/center/api/response"
	"apps/center/global"
	"apps/center/model"
	"apps/center/service"
	"apps/common/utils"
	"crypto/x509"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OnlineNode(c *gin.Context) {

	var nodes []model.ProbeInfo
	global.DB.Find(&nodes)
	c.JSON(200, response.Success(nodes))

}

func GenerateRSAKeyPair(c *gin.Context) {

	pub, pri, err := utils.GenerateRSAKeys(500)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(pri)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(pub)

	privateKeyBase64 := utils.EncodeKeyToBase64(privateKeyBytes)
	publicKeyBase64 := utils.EncodeKeyToBase64(publicKeyBytes)

	err = service.CENTER_INFO.SaveRSAKeyPair(publicKeyBase64, privateKeyBase64)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"publicKey": publicKeyBase64,
	})
}

func GetRSAKeyPairs(c *gin.Context) {
	pairs, err := service.CENTER_INFO.GetRSAKeyPairs()
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(pairs))
}
