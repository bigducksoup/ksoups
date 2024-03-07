package utils

import (
	"crypto/x509"
	"fmt"
	"testing"
)

func TestSliceToMap(t *testing.T) {

	type testStruct struct {
		Id   string
		Name string
	}

	s := []testStruct{
		{
			Id:   "1",
			Name: "1",
		},
		{
			Id:   "2",
			Name: "2",
		},
		{
			Id:   "3",
			Name: "3",
		},
	}

	m, err := SliceToMap[string, testStruct](s, "Id")

	if err != nil {
		t.Error(err)
	}

	if len(m) != 3 {
		t.Error("length not equal")
	}

	if m["1"].Name != "1" {
		t.Error("value not equal")
	}

	if m["2"].Name != "2" {
		t.Error("value not equal")
	}

	if m["3"].Name != "3" {
		t.Error("value not equal")
	}

}

func TestCrypto(t *testing.T) {

	pub, pri, err := GenerateRSAKeys(500)

	if err != nil {
		t.Fatal(err)
	}

	priBytes := x509.MarshalPKCS1PrivateKey(pri)
	privateKeyString := EncodeKeyToBase64(priBytes)
	fmt.Printf("privatekey is %s\n", privateKeyString)

	pubBytes := x509.MarshalPKCS1PublicKey(pub)
	publicKeyString := EncodeKeyToBase64(pubBytes)
	fmt.Printf("pubkey is %s\n", publicKeyString)

	decodedPriBytes, err := DecodeBase64ToKey(privateKeyString)
	if err != nil {
		t.Fatal(err)
	}

	decodedPubBytes, err := DecodeBase64ToKey(publicKeyString)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := ParsePrivateKey(decodedPriBytes)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ParsePublicKey(decodedPubBytes)
	if err != nil {
		t.Fatal(err)
	}

	originMessage := "hello"

	encryptedData, err := EncryptData([]byte(originMessage), publicKey)
	if err != nil {
		t.Fatal(err)
	}

	decryptedData, err := DecryptData(encryptedData, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(decryptedData))

}
