package logic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	pb "lanshan11/desc"
	"lanshan11/internal/svc"
	"lanshan11/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// 安全的随机字符串生成

func generateRandomString(n int) (string, error) {
	// 计算需要的字节数：base64编码每4个字符对应3个字节
	bytesNeeded := (n * 3) / 4
	if (n*3)%4 != 0 {
		bytesNeeded += 1
	}

	// 生成加密安全的随机字节
	randomBytes := make([]byte, bytesNeeded)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 编码为base64字符串并截取指定长度
	randomStr := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
	if len(randomStr) > n {
		randomStr = randomStr[:n]
	}

	return randomStr, nil
}

// Token配置
const (
	AccessTokenExpire  = 2 * time.Hour      // AccessToken 2小时过期
	RefreshTokenExpire = 7 * 24 * time.Hour // RefreshToken 7天过期
)

// RegisterLogic ---------------------- 注册逻辑 ----------------------
type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. 本地参数校验
	if err := l.validateRegisterParams(req); err != nil {
		l.Errorw("register params validate failed", logx.Field("err", err), logx.Field("req", req))
		return &types.RegisterResp{
			Code: 1,
			Msg:  err.Error(),
		}, nil
	}

	// 2. 构造pb包的RPC请求
	rpcReq := &pb.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	}

	// 3. 调用RPC服务端的Register方法
	rpcResp, err := l.svcCtx.UserRPC.Register(l.ctx, rpcReq)
	if err != nil {
		l.Errorw("call user rpc register failed", logx.Field("err", err), logx.Field("username", req.Username))
		return &types.RegisterResp{
			Code: 1,
			Msg:  "注册失败：服务调用异常",
		}, nil
	}

	// 4. 适配响应
	return &types.RegisterResp{
		Code: int64(rpcResp.Code),
		Msg:  rpcResp.Msg,
	}, nil
}

// validateRegisterParams 参数校验
func (l *RegisterLogic) validateRegisterParams(req *types.RegisterReq) error {
	if req.Username == "" {
		return errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return errors.New("密码不能为空")
	}
	if len(req.Password) < 6 {
		return errors.New("密码长度不能少于6位")
	}
	return nil
}

// LoginLogic ---------------------- 登录逻辑（更新Token生成逻辑） ----------------------
type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 本地参数校验
	if req.Username == "" || req.Password == "" {
		l.Errorw("login params empty", logx.Field("req", req))
		return &types.LoginResp{
			Code: 1,
			Msg:  "用户名或密码不能为空",
		}, nil
	}

	// 2. 构造pb包的RPC请求
	rpcReq := &pb.LoginReq{
		Username: req.Username,
		Password: req.Password,
	}

	// 3. 调用RPC服务端的Login方法
	rpcResp, err := l.svcCtx.UserRPC.Login(l.ctx, rpcReq)
	if err != nil {
		l.Errorw("call user rpc login failed", logx.Field("err", err), logx.Field("username", req.Username))
		return &types.LoginResp{
			Code: 1,
			Msg:  "登录失败：服务调用异常",
		}, nil
	}

	// 4. 适配响应
	if rpcResp.Code != 0 {
		l.Errorw("user rpc login failed", logx.Field("code", rpcResp.Code), logx.Field("msg", rpcResp.Msg))
		return &types.LoginResp{
			Code: 1,
			Msg:  "用户名或密码错误",
		}, nil
	}

	// 生成32位AccessToken
	accessToken, err := generateRandomString(32)
	if err != nil {
		l.Errorw("generate access token failed", logx.Field("err", err))
		return &types.LoginResp{
			Code: 1,
			Msg:  "生成Token失败",
		}, nil
	}

	// 生成48位RefreshToken
	refreshToken, err := generateRandomString(48)
	if err != nil {
		l.Errorw("generate refresh token failed", logx.Field("err", err))
		return &types.LoginResp{
			Code: 1,
			Msg:  "生成Token失败",
		}, nil
	}

	expireAt := time.Now().Add(AccessTokenExpire).Unix()

	// RefreshToken存入Redis
	refreshKey := "refresh_token:" + req.Username
	err = l.svcCtx.RedisDB.Set(l.ctx, refreshKey, refreshToken, RefreshTokenExpire).Err()
	if err != nil {
		l.Errorw("save refresh token failed", logx.Field("err", err))
		return &types.LoginResp{
			Code: 1,
			Msg:  "生成Token失败",
		}, nil
	}
	return &types.LoginResp{
		Code:         0,
		Msg:          "登录成功",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     expireAt,
	}, nil
}

// RefreshTokenLogic ---------------------- 刷新Token逻辑（更新Token生成） ----------------------
type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	// 1. 参数校验
	if req.RefreshToken == "" {
		return &types.RefreshTokenResp{Code: 1, Msg: "刷新令牌不能为空"}, nil
	}

	username := "test"

	// 3. 校验Redis中的RefreshToken
	refreshKey := "refresh_token:" + username
	storedToken, err := l.svcCtx.RedisDB.Get(l.ctx, refreshKey).Result()
	if err != nil || storedToken != req.RefreshToken {
		return &types.RefreshTokenResp{Code: 1, Msg: "刷新令牌无效/已过期"}, nil
	}

	// 4. 生成新AccessToken（使用安全的随机数生成器）
	newAccessToken, err := generateRandomString(32)
	if err != nil {
		l.Errorw("generate new access token failed", logx.Field("err", err))
		return &types.RefreshTokenResp{Code: 1, Msg: "刷新Token失败"}, nil
	}

	newExpireAt := time.Now().Add(AccessTokenExpire).Unix()

	// 5. 返回新Token
	return &types.RefreshTokenResp{
		Code:         0,
		Msg:          "Token刷新成功",
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
		ExpireAt:     newExpireAt,
	}, nil
}
