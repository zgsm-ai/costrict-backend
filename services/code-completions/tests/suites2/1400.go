package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"
)

// 简化版加密通信系统实现

// CryptoEngine 加密引擎
type CryptoEngine struct {
	key []byte
}

// NewCryptoEngine 创建新的加密引擎
func NewCryptoEngine(key string) (*CryptoEngine, error) {
	// 检查密钥长度
	if len(key) != 32 {
		return nil, fmt.Errorf("密钥长度必须为32字节")
	}

	return &CryptoEngine{
		key: []byte(key),
	}, nil
}

// Encrypt 加密数据
func (ce *CryptoEngine) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ce.key)
	if err != nil {
		return nil, fmt.Errorf("创建加密块失败: %v", err)
	}

	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM模式失败: %v", err)
	}

	// 创建随机数
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("生成随机数失败: %v", err)
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt 解密数据
func (ce *CryptoEngine) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ce.key)
	if err != nil {
		return nil, fmt.Errorf("创建加密块失败: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM模式失败: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("密文太短")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	return plaintext, nil
}

// Hash 哈希数据
func (ce *CryptoEngine) Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// SecureMessage 安全消息
type SecureMessage struct {
	ID      string    `json:"id"`
	Payload []byte    `json:"payload"`
	Hash    string    `json:"hash"`
	Time    time.Time `json:"time"`
}

// NewSecureMessage 创建新的安全消息
func NewSecureMessage(id string, payload []byte, engine *CryptoEngine) (*SecureMessage, error) {
	msg := &SecureMessage{
		ID:      id,
		Payload: payload,
		Hash:    engine.Hash(payload),
		Time:    time.Now(),
	}
	return msg, nil
}

// Verify 验证消息
func (sm *SecureMessage) Verify(engine *CryptoEngine) bool {
	expectedHash := engine.Hash(sm.Payload)
	return sm.Hash == expectedHash
}

<｜fim▁hole｜>

func main() {
	// 创建加密引擎
	engine, err := NewCryptoEngine("this-is-a-secret-key-32-bytes-long")
	if err != nil {
		log.Fatalf("创建加密引擎失败: %v", err)
	}

	fmt.Println("加密引擎创建成功")

	// 测试加密和解密
	plaintext := []byte("这是一条需要加密的测试消息")
	ciphertext, err := engine.Encrypt(plaintext)
	if err != nil {
		log.Printf("加密失败: %v", err)
	} else {
		fmt.Printf("加密成功，密文长度: %d\n", len(ciphertext))
	}

	// 测试解密
	decrypted, err := engine.Decrypt(ciphertext)
	if err != nil {
		log.Printf("解密失败: %v", err)
	} else {
		fmt.Printf("解密成功，明文: %s\n", string(decrypted))
	}

	// 测试哈希
	hash := engine.Hash(plaintext)
	fmt.Printf("哈希值: %s\n", hash)

	// 测试安全消息
	secureMsg, err := NewSecureMessage("msg-1", plaintext, engine)
	if err != nil {
		log.Printf("创建安全消息失败: %v", err)
	} else {
		// 验证消息
		if secureMsg.Verify(engine) {
			fmt.Println("安全消息验证成功")
		} else {
			fmt.Println("安全消息验证失败")
		}
	}

	// 等待一段时间
	time.Sleep(1 * time.Second)

	fmt.Println("简化版加密通信系统示例完成")
}
