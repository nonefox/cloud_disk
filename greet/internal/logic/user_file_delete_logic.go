package logic

import (
	"cloud_disk/greet/models"
	"context"
	"errors"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileDeleteLogic {
	return &UserFileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest, userIdentity string) (resp *types.UserFileDeleteResponse, err error) {
	//删除用户文件（这里的删除是软删除，在delete_at字段上面设置了删除时间）
	_, err = l.svcCtx.MysqlEngine.Where("user_identity = ? AND identity = ?", userIdentity, req.Identity).Delete(new(models.UserRepository))
	if err != nil {
		return nil, errors.New("删除用户文件失败")
	}
	return
}
