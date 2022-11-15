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

type UserFileNameChangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameChangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameChangeLogic {
	return &UserFileNameChangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameChangeLogic) UserFileNameChange(req *types.UserFileNameChangeRequest, userIdentity string) (resp *types.UserFileNameChangeResponse, err error) {
	//先判断在当前层级下是否有这个文件
	count, err := l.svcCtx.MysqlEngine.Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity = ?)",
		req.Name, req.Identity).Count(new(models.UserRepository))
	if err != nil {
		log.Printf("获取用户文件关联信息失败：%v", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("当前名称已被其他文件占用")
	}

	//通过前端传过来的文件名和文件Identity来更新文件名信息
	file := new(models.RepositoryPool)
	file.Name = req.Name
	_, err = l.svcCtx.MysqlEngine.Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Update(file)
	if err != nil {
		log.Printf("更新用户文件名失败：%v", err)
		return nil, err
	}

	return
}
