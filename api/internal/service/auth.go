package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qxsugar/bill/api/internal/model"
	"github.com/spf13/viper"
)

// AuthService 负责微信登录（code2session）与 JWT 签发/校验。
type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService: userService}
}

// claims 自定义 JWT 载荷，携带用户 id 与 openid。
type claims struct {
	UserId int64  `json:"uid"`
	Openid string `json:"openid"`
	jwt.RegisteredClaims
}

// code2sessionResp 微信 jscode2session 返回结构。
type code2sessionResp struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// LoginByCode 用微信 code 换 openid（配置缺失时走 dev 兜底），
// 找到或创建用户后签发 JWT，返回 token 与用户信息。
func (s *AuthService) LoginByCode(code string) (string, *model.User, error) {
	openid, err := s.code2openid(code)
	if err != nil {
		return "", nil, err
	}

	user, err := s.userService.GetOrCreateByOpenid(openid, randomNickname(), randomAvatar())
	if err != nil {
		return "", nil, err
	}

	token, err := s.signToken(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// code2openid 调用微信接口换取 openid；未配置 appid/secret 时返回 dev 兜底 openid。
func (s *AuthService) code2openid(code string) (string, error) {
	appid := viper.GetString("wechat.appid")
	secret := viper.GetString("wechat.secret")
	if appid == "" || secret == "" {
		// dev 兜底：同一个 code 稳定映射到同一个 openid，方便本地联调。
		return "dev_" + code, nil
	}

	endpoint := "https://api.weixin.qq.com/sns/jscode2session"
	query := url.Values{}
	query.Set("appid", appid)
	query.Set("secret", secret)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")

	cli := &http.Client{Timeout: 5 * time.Second}
	resp, err := cli.Get(endpoint + "?" + query.Encode())
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	var body code2sessionResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.ErrCode != 0 {
		return "", fmt.Errorf("wechat code2session failed: %d %s", body.ErrCode, body.ErrMsg)
	}
	if body.Openid == "" {
		return "", fmt.Errorf("wechat code2session returned empty openid")
	}
	return body.Openid, nil
}

// signToken 按配置的密钥与有效期签发 JWT。
func (s *AuthService) signToken(user *model.User) (string, error) {
	expireHours := viper.GetInt("jwt.expire_hours")
	now := time.Now()
	c := claims{
		UserId: user.Id,
		Openid: user.Openid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireHours) * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}

// ParseToken 校验 JWT 并返回其中的用户 id。
func (s *AuthService) ParseToken(tokenStr string) (int64, error) {
	c := &claims{}
	_, err := jwt.ParseWithClaims(tokenStr, c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return 0, err
	}
	return c.UserId, nil
}
