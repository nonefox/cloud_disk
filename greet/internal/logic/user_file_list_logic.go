package logic

import (
	"cloud_disk/greet/define"
	"cloud_disk/greet/models"
	"context"
	"log"
	"time"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest, userIdentity string) (resp *types.UserFileListResponse, err error) {
	//新建用户文件的存储结构和文件总数
	var ufList = make([]*types.UserFile, 0)

	//定义分页信息（）有参数就用参数值，没有就用默认值
	pageSize := req.Size
	if pageSize == 0 {
		pageSize = define.PageSize
	}
	pageIndex := req.Page
	if pageIndex == 0 {
		pageIndex = define.PageIndex
	}
	pageIndex = (pageIndex - 1) * pageSize //位移

	//查询并且返回ufList的数据，这里我们需要关联两张表一起查，把需要的数据组合成我们的userFile数据返回。
	err = l.svcCtx.MysqlEngine.Table("user_repository").Where("parent_id = ? AND user_identity = ?", req.Identity, userIdentity).
		Select("user_repository.id, user_repository.identity, user_repository.repository_identity, user_repository.ext,"+
			"user_repository.name, repository_pool.path, repository_pool.size").
		Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").
		Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format(define.Datetime)).
		Limit(pageSize, pageIndex).Find(&ufList)
	if err != nil {
		log.Printf("获取用户文件列表失败：%v", err)
		return
	}

	// 查询用户文件总数
	cnt, err := l.svcCtx.MysqlEngine.Where("parent_id = ? AND user_identity = ? ", req.Identity, userIdentity).Count(new(models.UserRepository))
	if err != nil {
		log.Printf("获取用户文件列表总数失败：%v", err)
		return
	}

	//定义响应类型
	resp = new(types.UserFileListResponse)
	resp.List = ufList
	resp.Count = cnt

	return resp, nil
}
