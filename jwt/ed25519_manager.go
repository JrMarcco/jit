package jwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-kit/bean/option"
	"github.com/golang-jwt/jwt/v5"
)

var _ Manager[any] = (*Ed25519Manager[any])(nil)

type Ed25519Manager[T any] struct {
	config ClaimsConfig

	priKey ed25519.PrivateKey // 私钥
	pubKey ed25519.PublicKey  // 公钥
}

func (m *Ed25519Manager[T]) Encrypt(data T) (string, error) {
	now := time.Now()
	cc := &CustomClaims[T]{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.Expiration)),
			ID:        m.config.JtiGenerator(),
		},
	}

	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, cc)
	return token.SignedString(m.priKey)
}

func (m *Ed25519Manager[T]) Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&CustomClaims[T]{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("[easy-kit] unexpected signing method: %v", token.Header["alg"])
			}
			return m.pubKey, nil
		},
		opts...,
	)
	if err != nil || !jwtToken.Valid {
		return CustomClaims[T]{}, fmt.Errorf("[easy-kit] failed to verify jwt token: %w", err)
	}
	cc, _ := jwtToken.Claims.(*CustomClaims[T])
	return *cc, nil
}

func WithClaimsConfig[T any](cfg ClaimsConfig) option.Opt[Ed25519Manager[T]] {
	return func(m *Ed25519Manager[T]) {
		m.config = cfg
	}
}

// NewEd25519Manager 创建使用 ed25519 加密算法的 jwt 管理器。
// 注意默认 token 过期时间为 24 小时。
func NewEd25519Manager[T any](priPem, pubPem string, opts ...option.Opt[Ed25519Manager[T]]) (*Ed25519Manager[T], error) {
	priKey, err := loadPrivateKey(priPem)
	if err != nil {
		return nil, err
	}
	pubKey, err := loadPublicKey(pubPem)
	if err != nil {
		return nil, err
	}

	m := &Ed25519Manager[T]{
		config: NewClaimsConfig(24 * time.Hour), // 默认 24 小时过期
		priKey: priKey,
		pubKey: pubKey,
	}
	option.Apply(m, opts...)
	return m, nil
}

// loadPrivateKey 加载私钥。
// 注意 PEM 块本身标注的是密钥对，而不是具体的 ed25519 密钥对。
// 所有标准公钥格式都需要先由 x509 包处理进行转换后类型断言才能获得 ed25519 密钥对。
func loadPrivateKey(priPem string) (ed25519.PrivateKey, error) {
	priKeyBlock, _ := pem.Decode([]byte(priPem))
	if priKeyBlock == nil {
		return nil, fmt.Errorf("[easy-kit] failed to decode private key PEM")
	}
	priKey, err := x509.ParsePKCS8PrivateKey(priKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("[easy-kit] failed to parse private key: %w", err)
	}
	return priKey.(ed25519.PrivateKey), nil
}

// loadPublicKey 加载公钥。
func loadPublicKey(pubPem string) (ed25519.PublicKey, error) {
	pubKeyBlock, _ := pem.Decode([]byte(pubPem))
	if pubKeyBlock == nil {
		return nil, fmt.Errorf("[easy-kit] failed to decode public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("[easy-kit] failed to parse public key: %w", err)
	}
	return publicKey.(ed25519.PublicKey), nil
}
