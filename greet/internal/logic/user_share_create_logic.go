package logic

import (
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"errors"
	"log"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserShareCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserShareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserShareCreateLogic {
	return &UserShareCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserShareCreateLogic) UserShareCreate(req *types.UserShareCreateRequest, userIdentity string) (resp *types.UserShareCreateResponse, err error) {
	ur := new(models.UserRepository)
	//先通过用户文件关联表identity来判断用户是否有该文件
	get, err := l.svcCtx.MysqlEngine.Where("identity = ? AND user_identity = ?", req.UserRepositoryIdentity, userIdentity).Get(ur)
	if err != nil {
		log.Printf("用户分享查询用户文件失败：%v", err)
		return nil, err
	}
	if !get {
		return nil, errors.New("用户没有该文件，无法分享")
	}

	//存在我就插入分享记录到share_basic表中
	shareData := models.ShareBasic{
		Identity:               tools.GetUUID(),
		UserIdentity:           userIdentity,
		UserRepositoryIdentity: req.UserRepositoryIdentity, //用户文件关联Identity
		RepositoryIdentity:     ur.RepositoryIdentity,      //文件Identity
		ExpiredTime:            req.ExpiredTime,            //过期时间
	}
	//插入数据
	_, err = l.svcCtx.MysqlEngine.InsertOne(shareData)
	if err != nil {
		log.Printf("插入用户分享记录失败：%v", err)
		return nil, err
	}
	//定义响应类型
	resp = &types.UserShareCreateResponse{
		Identity: ur.RepositoryIdentity,
	}
	return resp, nil
}
