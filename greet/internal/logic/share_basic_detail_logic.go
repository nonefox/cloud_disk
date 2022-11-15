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

type ShareBasicDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.ShareBasicDetailRequest) (resp *types.ShareBasicDetailResponse, err error) {
	//确认该identity的资源分享是否存在
	share := new(models.ShareBasic)
	get, err := l.svcCtx.MysqlEngine.Where("identity = ?", req.Identity).Get(share)
	if err != nil {
		log.Printf("查询分享资源失败：%v", err)
		return nil, err
	}
	if !get {
		return nil, errors.New("当前分享资源不存在")
	}
	//更新分享资源的点击数(+1)
	share.ClickNum = share.ClickNum + 1
	_, err = l.svcCtx.MysqlEngine.Update(share)
	if err != nil {
		log.Printf("更新分享资源点击数失败：%v", err)
		return nil, err
	}

	//我们通过连表查询出需要返回给用户的分享资源详情信息
	resp = new(types.ShareBasicDetailResponse)
	_, err = l.svcCtx.MysqlEngine.Table("share_basic").
		Select("share_basic.repository_identity, user_repository.name, repository_pool.ext, repository_pool.size, repository_pool.path").
		Join("LEFT", "repository_pool", "share_basic.repository_identity = repository_pool.identity").
		Join("LEFT", "user_repository", "user_repository.identity = share_basic.user_repository_identity").
		Where("share_basic.identity = ?", req.Identity).Get(resp)

	if err != nil {
		return nil, errors.New("分享资源详情获取失败")
	}
	return resp, nil
}
