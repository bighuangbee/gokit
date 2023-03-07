/**
 * @desc //TODO $
 * @param $
 * @return $
 **/
package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	text := "1627992074_0848a8629da2049a6ea78374e8dc77bb_2_2021-11-09 15:04:09_00000000029"
	AesKey := []byte("d1a5f9c5877bda90dc1d724f6173bf04") //秘钥长度为16的倍数
	fmt.Printf("明文: %s\n秘钥: %s\n", text, string(AesKey))
	encrypted, err := AesEncrypt([]byte(text), AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))
	origin, err := AesDecrypt(encrypted, AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("解密后明文: %s\n", string(origin))
}

func TestAESstr(t *testing.T) {
	text := "d1a5f9c5877bda90dc1d724f6173bf04"
	AesKey := "6d1a5f9c577bda90dc1d724f6173bf04" //秘钥长度为16的倍数
	fmt.Printf("明文: %s\n秘钥: %s\n", text, AesKey,)
	encrypted, err := AesEncryptStr(text, AesKey)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("密文:", encrypted)

	origin, err := AesDecryptStr(encrypted, AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("解密后明文: %s\n", string(origin))
}


func TestAESLicense(t *testing.T) {
	AesKey := "111888" //秘钥长度为16的倍数
	encrypted, err := AesDecryptStr("r1hWlBPteGkfR73GgA5euUqbkN7i6TzmpwhDR-4cSFn2K6XhRPxQ1MTEL5f7vy5IppYyzaaQAVODW5Nmp_kyzhl-2j9Y8xuoUiBm-ti87NXiGD9HEtbS3udp7-Rc09o=", AesKey)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("密文:", encrypted)

	origin, err := AesDecryptStr(encrypted, AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("解密后明文: %s\n", string(origin))
}
