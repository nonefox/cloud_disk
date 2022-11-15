package logic

import (
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"errors"
	"log"
	"time"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCodeLogic) SendCode(req *types.SendCodeRequest) (resp *types.SendCodeResponse, err error) {
	user := new(models.UserBasic)
	//判断该邮箱是否被注册
	count, err := l.svcCtx.MysqlEngine.Where("email = ?", req.Email).Count(user)
	if err != nil {
		log.Printf("通过邮箱获取数据失败：%v", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该邮箱已注册")
	}
	//发送邮箱验证码
	code := tools.GenerateCode() //生成验证码
	//存储验证码（放在redis中,以email作为key，code作为值）
	l.svcCtx.RedDB.Set(l.ctx, req.Email, code, time.Second*300)
	err = tools.SendCode(req.Email, code) //发送验证码
	if err != nil {
		return nil, errors.New("发送验证码失败")
	}
	return nil, nil
}
