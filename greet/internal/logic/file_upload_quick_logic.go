package logic

import (
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"log"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadQuickLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadQuickLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadQuickLogic {
	return &FileUploadQuickLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FileUploadQuick 文件快传，当用户上传的文件已经存在于中心存储池中时（通过文件md5值来比较），我们就直接快传，然后更新用户文件关联表
func (l *FileUploadQuickLogic) FileUploadQuick(req *types.FileUploadQuickRequest, userIdentity string) (resp *types.FileUploadQuickResponse, err error) {
	//先判断该文件是否存在中心存储池中
	rp := new(models.RepositoryPool)
	get, err := l.svcCtx.MysqlEngine.Where("hash = ?", req.Md5).Get(rp)
	if err != nil {
		log.Printf("文件查询失败：%v", err)
		return nil, err
	}

	//定义响应类型
	resp = new(types.FileUploadQuickResponse)
	if get { //如果存在直接快传更新用户文件关联表
		//把新的用户存储关系，加入到user_repository表中
		ur := &models.UserRepository{
			Identity:           tools.GetUUID(),
			UserIdentity:       userIdentity,
			ParentId:           req.ParentId,
			RepositoryIdentity: rp.Identity,
			Ext:                rp.Ext,
			Name:               rp.Name,
		}
		_, err = l.svcCtx.MysqlEngine.InsertOne(ur)
		if err != nil {
			log.Printf("插入用户存储关系表失败：%v", err)
			return nil, err
		}
		resp.Identity = rp.Identity //把资源identity返回给前端
		return resp, nil
	} else {
		//如果没有我们就分片上传
		key, UploadID, err := tools.InitTenCosPart(rp.Ext)
		if err != nil {
			log.Printf("初始化分片上传失败：%v", err)
			return nil, err
		}
		//把分片传输的key和UploadID也一起传给前端
		resp.Key = key
		resp.UploadId = UploadID
	}
	return resp, nil
}
