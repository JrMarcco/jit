package xjwt

import (
	"time"

	"github.com/JrMarcco/easy-kit/bean/option"
	"github.com/golang-jwt/jwt/v5"
)

// Manager jwt 管理器抽象
type Manager[T any] interface {
	Encrypt(data T) (string, error)
	Decrypt(token string, opts ...jwt.ParserOption) (CustomClaims[T], error)
}

type CustomClaims[T any] struct {
	jwt.RegisteredClaims
	Data T
}

// ClaimsConfig jwt claims 扩展配置项
type ClaimsConfig struct {
	Issuer       string        // 签发人
	Expiration   time.Duration // 有效期
	JtiGenerator func() string // jwt id 生成方法
}

func WithIssuer(issuer string) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.Issuer = issuer
	}
}

func WithJtiGenerator(jtiGenerator func() string) option.Opt[ClaimsConfig] {
	return func(cfg *ClaimsConfig) {
		cfg.JtiGenerator = jtiGenerator
	}
}

// NewClaimsConfig 创建自定义 claims 配置信息。
// 过期时间必须显示指定。
// 其他参数可以通过 option.Opt 自定义。
func NewClaimsConfig(expiration time.Duration, opts ...option.Opt[ClaimsConfig]) ClaimsConfig {
	cfg := ClaimsConfig{
		Issuer:       "easy-kit", // 默认签发人
		Expiration:   expiration,
		JtiGenerator: func() string { return "" },
	}

	option.Apply(&cfg, opts...)

	return cfg
}
