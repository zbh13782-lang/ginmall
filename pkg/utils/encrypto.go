package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type Encryption struct {
	key string
}

var Encrypt *Encryption

func init() {
	Encrypt = NewEncryption()
}

func NewEncryption() *Encryption {
	return &Encryption{}
}

// PKCS7 填充（标准填充方式，替代原有简易填充）
// blockSize：AES块大小固定为16字节

func PadPwd(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtex := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtex...)
}

func UnPadPwd(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("长度为0")
	}
	padding := int(src[length-1])
	if padding > length || padding == 0 {
		return nil, errors.New("长度不合法")
	}
	for i := 0; i < padding; i++ {
		if src[length-padding+i] != byte(padding) {
			return nil, errors.New("填充内容不合法")
		}
	}
	// 去除填充部分
	return src[:length-padding], nil
}

func (k *Encryption) SetKey(key string) {
	k.key = key
}
//加密
func (k *Encryption) AesEncoding(src string) (string, error) {
	// 1. 转换密钥和明文为字节数组
	keyBytes := []byte(k.key)
	srcBytes := []byte(src)

	// 2. 校验密钥长度
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", errors.New("密钥不合法（长度必须为16/24/32字节）：" + err.Error())
	}

	// 3. PKCS7填充明文（AES块大小固定为16字节）
	blockSize := block.BlockSize()
	paddedSrc := PadPwd(srcBytes, blockSize)

	// 4. 生成初始化向量IV（CBC模式必须使用，长度等于AES块大小16字节）
	// 注意：生产环境中IV应随机生成并与密文一起存储/传输，此处为方便演示，可根据业务需求调整
	// 推荐方案：IV随机生成（16字节），加密后密文 = Base64(IV + 加密数据)
	iv := bytes.Repeat([]byte("1234567890abcdef"), 1) // 演示用固定IV，生产环境请替换为随机IV
	if len(iv) != blockSize {
		return "", errors.New("IV长度必须等于AES块大小（16字节）")
	}

	// 5. 初始化CBC模式加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 6. 执行加密（CBC模式要求明文长度是块大小的整数倍，已通过填充保证）
	crypted := make([]byte, len(paddedSrc))
	mode.CryptBlocks(crypted, paddedSrc)

	// 7. Base64编码密文，方便存储和传输
	return base64.StdEncoding.EncodeToString(crypted), nil
}

//解密
func (k *Encryption) AesDecoding(pwd string) (string, error) {
	// 1. Base64解码密文
	cryptedBytes, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return "", errors.New("密文Base64解码失败：" + err.Error())
	}

	// 2. 转换密钥为字节数组并校验
	keyBytes := []byte(k.key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", errors.New("密钥不合法（长度必须为16/24/32字节）：" + err.Error())
	}

	// 3. 校验密文长度（必须是块大小的整数倍）
	blockSize := block.BlockSize()
	if len(cryptedBytes)%blockSize != 0 {
		return "", errors.New("密文长度不合法（必须是16字节的整数倍）")
	}

	// 4. 初始化向量IV（需与加密时使用的IV一致）
	iv := bytes.Repeat([]byte("1234567890abcdef"), 1)
	if len(iv) != blockSize {
		return "", errors.New("IV长度必须等于AES块大小（16字节）")
	}

	// 5. 初始化CBC模式解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 6. 执行解密
	decrypted := make([]byte, len(cryptedBytes))
	mode.CryptBlocks(decrypted, cryptedBytes)

	// 7. PKCS7去填充，得到原始明文
	origSrc, err := UnPadPwd(decrypted)
	if err != nil {
		return "", errors.New("去填充失败：" + err.Error())
	}

	// 8. 转换为字符串并返回
	return string(origSrc), nil
}