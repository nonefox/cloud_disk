package logic

import (
	"cloud_disk/greet/define"
	"cloud_disk/greet/tools"
	"context"
	"log"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

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

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest, authorization string) (resp *types.RefreshTokenResponse, err error) {
	//获取前端用户认证中的用户声明信息
	userClaim, err := tools.AnalyseToken(authorization)
	if err != nil {
		log.Printf("解析用户认证authorization失败：%v", err)
		return nil, err
	}

	//通过用户声明中的信息去生成我们新的token和refreshToken
	token, err := tools.GenerateToken(userClaim.Id, userClaim.Identity, userClaim.Name, int64(define.TokenTime))
	if err != nil {
		log.Printf("用户token生成失败：%v", err)
		return nil, err
	}
	refreshToken, err := tools.GenerateToken(userClaim.Id, userClaim.Identity, userClaim.Name, int64(define.RefreshTokenTime))
	if err != nil {
		log.Printf("用户token生成失败：%v", err)
		return nil, err
	}

	//生成响应结构
	resp = new(types.RefreshTokenResponse)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return resp, nil
}
