package connect

import (
	"apps/common/message/data"
	"apps/common/utils"
)

func GenerateRegisterInfo(pubKey string, willBeEncrypted string) (data.RegisterInfo, error) {
	bytes, err := utils.DecodeBase64ToKey(pubKey)

	if err != nil {
		return data.RegisterInfo{}, err
	}

	publicKey, err := utils.ParsePublicKey(bytes)

	if err != nil {
		return data.RegisterInfo{}, err
	}

	encryptData, err := utils.EncryptData([]byte(willBeEncrypted), publicKey)

	if err != nil {
		return data.RegisterInfo{}, err
	}

	md5 := utils.Md5([]byte(pubKey))

	return data.RegisterInfo{
		Name:          willBeEncrypted,
		EncryptedName: encryptData,
		PublicKeyMd5:  md5,
	}, nil
}
