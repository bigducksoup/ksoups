package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

func GenerateRSAKeys(keySize int) (pub *rsa.PublicKey, pri *rsa.PrivateKey, err error) {

	pri, err = rsa.GenerateKey(rand.Reader, keySize)

	if err != nil {
		return nil, nil, err
	}

	return &pri.PublicKey, pri, nil
}

func EncodeKeyToBase64(keyBytes []byte) string {
	return base64.StdEncoding.EncodeToString(keyBytes)
}

func DecodeBase64ToKey(keyString string) ([]byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

func ParsePrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(privateKeyBytes)
}

func ParsePublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	return x509.ParsePKCS1PublicKey(publicKeyBytes)
}

func EncryptData(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

func DecryptData(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}
