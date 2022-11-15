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

type UserFolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderCreateLogic {
	return &UserFolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateRequest, userIdentity string) (resp *types.UserFolderCreateResponse, err error) {
	//先判断该文件夹是否存在用户存储关系表中
	count, err := l.svcCtx.MysqlEngine.Where("name = ? AND parent_id = ?", req.Name, req.ParentId).Count(new(models.UserRepository))
	if err != nil {
		log.Printf("查询文件夹信息失败：%v", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该名称文件夹已存在当前层级")
	}

	//如果没有我们创建新的文件夹
	userFolder := new(models.UserRepository)
	userFolder.Identity = tools.GetUUID()
	userFolder.UserIdentity = userIdentity
	userFolder.Name = req.Name
	userFolder.ParentId = req.ParentId
	_, err = l.svcCtx.MysqlEngine.InsertOne(userFolder)
	if err != nil {
		log.Printf("创建新文件夹失败：%v", err)
		return nil, err
	}

	//返回响应类型
	return &types.UserFolderCreateResponse{
		Identity: userFolder.Identity,
	}, nil
}
