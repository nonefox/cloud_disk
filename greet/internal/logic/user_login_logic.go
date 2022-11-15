package logic

import (
	"cloud_disk/greet/define"
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"errors"
	"log"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	//定义数据表到数据结构的映射
	user := new(models.UserBasic)
	//定义响应类型
	res := new(types.LoginResponse)
	//首先查询用户是否存在(通过用户名和密码)
	result, err := l.svcCtx.MysqlEngine.Where("name = ? AND password = ?", req.Name, tools.GetMd5(req.Password)).Get(user)
	if err != nil {
		log.Printf("用户登陆失败：%v", err)
		return nil, err
	}
	if !result {
		log.Printf("用户信息不存在")
		return nil, errors.New("用户名或密码错误")
	}

	//然后生成当前用户的token信息
	token, err := tools.GenerateToken(user.Id, user.Identity, user.Name, int64(define.TokenTime))
	if err != nil {
		log.Printf("生成用户Token失败：%v", err)
		return nil, err
	}

	//生成刷新token信息
	refreshToken, err := tools.GenerateToken(user.Id, user.Identity, user.Name, int64(define.RefreshTokenTime))
	if err != nil {
		log.Printf("生成用户refreshToken失败：%v", err)
		return nil, err
	}

	res.Token = token
	res.RefreshToken = refreshToken

	return res, nil
}
