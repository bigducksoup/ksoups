package info

import (
	"apps/center/model"
	"apps/common/utils"
	"gorm.io/gorm"
)

type CenterInfoService struct {
	Db *gorm.DB
}

func (c *CenterInfoService) SaveRSAKeyPair(publicKey string, privateKey string) error {

	md5 := utils.Md5([]byte(publicKey))

	return c.Db.Save(&model.RegisterKey{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		PublicKeyMd5: md5,
	}).Error
}

func (c *CenterInfoService) GetRSAKeyPairs() (pairs []model.RegisterKey, err error) {
	err = c.Db.Find(&pairs).Error
	if err != nil {
		return nil, err
	}
	return
}

func (c *CenterInfoService) GetKeyPairByPublicKeyMd5(md5 string) (*model.RegisterKey, error) {

	key := model.RegisterKey{}
	err := c.Db.Where("public_key_md5 = ?", md5).Find(&key).Error

	if err != nil {
		return nil, err
	}

	return &key, nil
}
