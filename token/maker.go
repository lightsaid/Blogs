package token

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var ErrInvalidToken = errors.New("token 无效")

// Maker 定义 jwt 接口两个核心接口
type Maker interface {
	// GenToken 根据用户ID生成有时效的token
	GenToken(*Payload) (string, error)

	// ParseToken 解析并验证token是否有效
	ParseToken(token string) (*Payload, error)
}

// Payload 定义Token负载数据
type Payload struct {
	ID        uuid.UUID `json:"id"`        // Token 唯一标识
	UserID    int64     `json:"userId"`    // 用户ID，Token持有者/使用者
	IssuedAt  time.Time `json:"issuedAt"`  // 签发时间
	ExpiredAt time.Time `json:"expiredAt"` // 过期时间
}

// NewPayload 创建一个 Payload 提供给 GenToken 方法使用
func NewPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	payload.ID = tokenID

	return payload, nil
}
