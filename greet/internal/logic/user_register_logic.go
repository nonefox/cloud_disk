package logic

import (
	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"errors"
	"log"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {
	//邮箱判断已经在sendcode中做了，这里就直接获取验证码（存在redis中,key是email）
	code, err := l.svcCtx.RedDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		log.Printf("redis获取验证码失败：%v", err)
		return nil, err
	}
	if code != req.Code {
		return nil, errors.New("验证码错误")
	}

	user := new(models.UserBasic)

	//判断用户是否已经存在
	count, err := l.svcCtx.MysqlEngine.Where("name = ?", req.Name).Count(user)
	if err != nil {
		log.Printf("判断用户是否已经存在失败：%v", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该用户已经注册")
	}

	//最后我们把新用户的数据存入表中
	user.Identity = tools.GetUUID()
	user.Name = req.Name
	user.Email = req.Email
	req.Password = tools.GetMd5(req.Password)

	_, err = l.svcCtx.MysqlEngine.InsertOne(user)
	if err != nil {
		log.Printf("保存新用户数据失败：%v", err)
		return nil, err
	}

	return nil, nil
}
