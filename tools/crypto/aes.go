/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func  AesEncryptStr(origData, key string)(string, error) {
	crypto, err :=  AesEncrypt([]byte(origData), []byte(key))
	if err != nil{
		return "", err
	}
	return base64.URLEncoding.EncodeToString(crypto), err
}


//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func AesDecryptStr(crypteStr, key string)(string, error) {
	crypted, err := base64.URLEncoding.DecodeString(crypteStr)
	if err != nil{
		return "", nil
	}
	crypto, err :=  AesDecrypt(crypted, []byte(key))
	return string(crypto), err
}

func AesCbcDecrypt(crypted, key, iv []byte)([]byte, error){
	block , err := aes.NewCipher(key)
	if err != nil{
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(crypted) < blockSize || len(crypted) % blockSize != 0{
		return nil, errors.New("crypted size err")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(crypted, crypted)

	data := PKCS7UnPadding(crypted)
	return data, nil
}
