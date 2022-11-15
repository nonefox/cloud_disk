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

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest, userIdentity string) (resp *types.UserFileMoveResponse, err error) {
	//要更换文件夹的位置，需要有目标文件夹的Identity（存放在user_repository表中）
	parent := new(models.UserRepository)
	//依据目标文件夹的Identity来获取对应ID
	get, err := l.svcCtx.MysqlEngine.Where("identity = ? AND user_identity = ?", req.ParentIdnetity, userIdentity).Get(parent)
	if err != nil {
		log.Printf("获取目标文件夹ID失败：%v", err)
		return nil, err
	}
	if !get {
		return nil, errors.New("目标文件夹不存在")
	}

	//更新对应文件的parent_id(代表文件层级)
	l.svcCtx.MysqlEngine.Where("identity = ?", req.Idnetity).Update(models.UserRepository{
		ParentId: int64(parent.Id),
	})

	return
}
