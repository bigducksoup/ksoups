package info

import (
	"apps/center/model"
	"apps/common/utils"
	"gorm.io/gorm"
)

type CenterInfoService struct {
	Db *gorm.DB
}

func (c *CenterInfoService) SaveRSAKeyPair(publicKey string, privateKey string, name string) error {

	md5 := utils.Md5([]byte(publicKey))

	return c.Db.Save(&model.RegisterKey{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		PublicKeyMd5: md5,
		Name:         name,
	}).Error
}

func (c *CenterInfoService) DeleteRSAKeyPair(id uint) error {
	return c.Db.Delete(&model.RegisterKey{}, id).Error
}

func (c *CenterInfoService) GetRSAKeyPairs() (pairs []model.RegisterKey, err error) {

	var vo []struct {
		model.BaseModel
		PublicKey    string `json:"public_key"`
		PublicKeyMd5 string `json:"public_key_md5"`
		Name         string `json:"name"`
	}

	err = c.Db.Model(&model.RegisterKey{}).Find(&vo).Error
	if err != nil {
		return nil, err
	}

	for _, v := range vo {
		pairs = append(pairs, model.RegisterKey{
			BaseModel:    v.BaseModel,
			PublicKey:    v.PublicKey,
			PublicKeyMd5: v.PublicKeyMd5,
			Name:         v.Name,
		})
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
