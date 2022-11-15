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

type UserShareSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserShareSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserShareSaveLogic {
	return &UserShareSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserShareSave 把其他用户分享的资源，保存到自己的用户文件关联表中
func (l *UserShareSaveLogic) UserShareSave(req *types.UserShareSaveRequest, userIdentity string) (resp *types.UserShareSaveResponse, err error) {
	//通过前端传过来的资源的identity，判断资源是否存在
	repPool := new(models.RepositoryPool)
	get, err := l.svcCtx.MysqlEngine.Where("identity = ?", req.RepositoryIdentity).Get(repPool)
	if err != nil {
		log.Printf("保存用户分享查询数据失败：%v", err)
		return nil, err
	}
	if !get {
		return nil, errors.New("该资源不存在")
	}

	//把资源信息保存到用户文件关联表中
	userRep := &models.UserRepository{
		Identity:           tools.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                repPool.Ext,
		Name:               repPool.Name,
	}
	_, err = l.svcCtx.MysqlEngine.InsertOne(userRep)
	if err != nil {
		log.Printf("用户保存资源失败：%v", err)
		return nil, err
	}

	//定义响应类型
	resp = new(types.UserShareSaveResponse)
	resp.Identity = userRep.Identity
	return resp, nil
}
