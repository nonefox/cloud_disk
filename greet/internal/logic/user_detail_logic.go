package logic

import (
	"cloud_disk/greet/models"
	"context"
	"errors"
	"log"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailRequest) (resp *types.UserDetailResponse, err error) {
	//定义数据表到数据结构的映射
	user := new(models.UserBasic)
	//定义响应类型
	res := new(types.UserDetailResponse)

	//通过用户的Identity来获取用户的detail信息
	result, err := l.svcCtx.MysqlEngine.Where("identity = ?", req.Identity).Get(user)
	if err != nil {
		log.Printf("通过identity获取用户detail失败：%v", err)
		return nil, err
	}
	if !result {
		log.Printf("当前用户数据未找到")
		return nil, errors.New("用户不存在")
	}
	res.Name = user.Name
	res.Email = user.Email
	return res, nil

}
