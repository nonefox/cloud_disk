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

type UserRepositoryRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRepositoryRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRepositoryRelationLogic {
	return &UserRepositoryRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserRepositoryRelation 需要用到当前用户的Identity，把放在请求的header中的UserIdentity传过来（在auth_middleware中放进去的）
func (l *UserRepositoryRelationLogic) UserRepositoryRelation(req *types.UserRepositoryRelationRequest, userIdentity string) (resp *types.UserRepositoryRelationResponse, err error) {
	//判断该资源是否存在repository_pool表中
	get, err := l.svcCtx.MysqlEngine.Where("identity = ?", req.RepositoryIdentity).Get(new(models.RepositoryPool))
	if err != nil {
		return nil, err
	}
	if !get {
		return nil, errors.New("该资源不存在")
	}

	//把新的用户存储关系，加入到user_repository表中
	ur := &models.UserRepository{
		Identity:           tools.GetUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                req.Ext,
		Name:               req.Name,
	}
	_, err = l.svcCtx.MysqlEngine.InsertOne(ur)
	if err != nil {
		log.Printf("插入用户存储关系表失败：%v", err)
		return nil, err
	}
	return
}
