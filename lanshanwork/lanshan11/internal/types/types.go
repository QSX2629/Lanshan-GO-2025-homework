package types

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// RegisterResp 注册响应
type RegisterResp struct {
	Code int64  `json:"code"` // 状态码：0成功，1失败
	Msg  string `json:"msg"`  // 提示信息
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// LoginResp 登录响应（补充Token字段）
type LoginResp struct {
	Code         int64  `json:"code"`          // 状态码
	Msg          string `json:"msg"`           // 提示信息
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpireAt     int64  `json:"expire_at"`     // AccessToken过期时间（秒级时间戳）
}

// RefreshTokenReq 刷新Token请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}

// RefreshTokenResp 刷新Token响应
type RefreshTokenResp struct {
	Code         int64  `json:"code"`          // 状态码
	Msg          string `json:"msg"`           // 提示信息
	AccessToken  string `json:"access_token"`  // 新访问令牌
	RefreshToken string `json:"refresh_token"` // 新刷新令牌（可选）
	ExpireAt     int64  `json:"expire_at"`     // 新Token过期时间
}
