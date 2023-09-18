package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// NOTE:
// 本包参考 Alexedwrads 大佬博客文章：
// https://www.alexedwards.net/blog/working-with-cookies-in-go

var (
	ErrValueTooLong = errors.New("cookie值太长了")
	ErrInvalidValue = errors.New("无效的cookie值")
)

// 每个Cookie最大值4k
const cookieMaxSize = 4096

// Write 简单地写入Cookie，仅做URL编码
func Write(w http.ResponseWriter, cookie http.Cookie) error {
	// 将 cookie 值编码
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > cookieMaxSize {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)

	return nil
}

// Read 简单地读取cookie，仅做URL解码
func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

// WriteSigned 先将cookie使用sha256算法签名再写入cookie，防止篡改，并没加密
func WriteSigned(w http.ResponseWriter, cookie http.Cookie, secretKey string) error {

	// NOTE:
	// HMAC-SHA256是一种算法，它是基于哈希函数SHA-256（Secure Hash Algorithm 256-bit）的消息认证码算法。
	// 并且生成的hash值时固定32个字节，即可 256 = 32 × 8；
	// 再者HMAC-SHA256单向哈希函数，不可逆，因此不能从签名恢复原始消息。
	// 签名的目的是验证消息的完整性和认证发送者，而不是用于加密消息。
	// 相同的消息经过HMAC-SHA256得到一致的hash

	// 因此实现防篡改签名的大致原理是：
	// 首先将消息使用HMAC-SHA225算法+密钥生成签名，
	// 然后获取消息后解开消息内容，使用同一个密钥再次生成签名，对比两次签名是否一致。

	// 具体实施：
	// 前提写cookie必须知道name、value，而读cookie仅知道name即可。
	// 1. 写的时候，将： （name + value） -> 签名 = signature， 签名的长度是固定的32个字节
	// 2. 设置新的 cookie.value = signature + value
	// 3. 写入cookie

	// 4. cookie通过name 读取到 value，此时的 value = signature + value
	// 5. 因此分别截取到原来的签名和真实的value
	// 6. 在此利用name和真实value重新生成签名，两次签名对比，一致就没有被篡改

	// 使用secretKey密钥创建一个HMAC-SHA256实例
	mac := hmac.New(sha256.New, []byte(secretKey))

	// 将要签名的信息（name、value）通过Write写入mac
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))

	// 获取最终的签名结果值
	signature := mac.Sum(nil)

	// 重新设置 value = 签名  + value
	cookie.Value = string(signature) + cookie.Value

	// 写入 http cookie
	return Write(w, cookie)
}

// ReadSigned
func ReadSigned(r *http.Request, name string, secretKey string) (string, error) {
	// 通过name获取带有签名 + value 的 signedValue
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}
	// 截取签名
	signature := signedValue[:sha256.Size]
	// 获取value
	value := signedValue[sha256.Size:]
	// 使用密钥重新生成一个hash
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	// 得到新的签名
	expectedSignature := mac.Sum(nil)
	// 对比两个签名
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return value, nil
}

// WriteEncrypted 使用对称加密后在写入cookie
func WriteEncrypted(w http.ResponseWriter, cookie http.Cookie, secretKey string) error {
	// NOTE：
	// 在 AES 中，密钥长度可以是 16、24 或 32 字节，分别对应 AES-128、AES-192 和 AES-256。选择一个合适的密钥长度，并确保密钥的安全性。
	// AES 加密算法要求输入数据的长度必须是块长度的整数倍（块长度为 16 字节）

	// 大概步骤：
	// 1. 选择密钥，16、24、32字节，创建aes加密块
	// 2. 填充数据
	// 3. 加密 Encrypt
	// 4. 解密 Decrypt

	// 创建一个AES加密块， key 可以时 16、24、32字节
	// 下面使用的cipher.NewGCM AES-128加密算法，因此密钥必须是16位
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return err
	}

	// 创建一个 AES-GCM 实例，AES-GCM 实现了 AEAD（Authenticated Encryption with Associated Data），
	// 同时进行加密和身份验证，这样可以确保加密数据的完整性。
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 返回 AES-GCM 加密中使用的 nonce（随机数）的大小
	nonce := make([]byte, aesGCM.NonceSize())
	// 从加密库的随机数生成器中生成随机的 nonce 值。nonce 是一个只使用一次的值，用于确保相同明文在不同的加密过程中生成不同的密文。
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	// 使用 AES-GCM 进行加密
	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = string(encryptedValue)

	return Write(w, cookie)
}

// ReadEncrypted 从cookie 读取value并解密返回
func ReadEncrypted(r *http.Request, name string, secretKey string) (string, error) {
	encryptedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptedValue) < nonceSize {
		return "", ErrInvalidValue
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", ErrInvalidValue
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", ErrInvalidValue
	}

	if expectedName != name {
		return "", ErrInvalidValue
	}

	return value, nil
}

// NOTE: 简单介绍一下算法
/*
在 Go 语言中，加密和解密是通过标准库中的 crypto 包来实现的。
crypto 包提供了多种加密算法和密码学功能，包括对称加密、非对称加密、哈希函数、数字签名等。

对称加密：
对称加密使用相同的密钥（称为对称密钥）来加密和解密数据。
常用的对称加密算法有 AES (Advanced Encryption Standard) 和 DES (Data Encryption Standard)。
对称加密适用于保护数据的机密性，但在数据传输和存储过程中需要安全地共享密钥。

非对称加密：
非对称加密使用一对密钥，一个公钥和一个私钥，来加密和解密数据。
公钥可以用于加密数据，但只有持有相应私钥的接收方才能解密数据。
常用的非对称加密算法有 RSA (Rivest–Shamir–Adleman) 和 ECC (Elliptic Curve Cryptography)。
非对称加密用于实现加密通信和数字签名。

哈希函数：
哈希函数将任意长度的数据映射为固定长度的哈希值。
哈希函数是单向的，即从哈希值无法恢复原始数据。
常用的哈希函数有 SHA-256 和 SHA-512。哈希函数用于数据完整性校验和密码散列。

数字签名：
数字签名结合了哈希函数和非对称加密，用于验证数据的来源和完整性。
数据的发送方使用私钥对数据的哈希值进行签名，接收方使用相应公钥验证签名的有效性。

在 Go 中，crypto/aes 包提供了对 AES 对称加密算法的支持，
crypto/rsa 包提供了对 RSA 非对称加密算法的支持，crypto/sha256 包提供了 SHA-256 哈希函数的支持等。

*/
