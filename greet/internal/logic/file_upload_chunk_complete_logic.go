package logic

import (
	"cloud_disk/greet/define"
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"

	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadChunkCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkCompleteLogic {
	return &FileUploadChunkCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkCompleteLogic) FileUploadChunkComplete(req *types.FileUploadChunkCompleteRequest, userIdentity string) (resp *types.FileUploadChunkCompleteResponse, err error) {
	//首先把前端传过来的object分片验证信息，重新分装成为cos的object类型
	object := make([]cos.Object, 0)
	for _, v := range req.CosObjects {
		object = append(object, cos.Object{
			ETag:       v.Etag,
			PartNumber: v.PartNumber,
		})
	}

	//调用分片完成的函数
	err = tools.TenCosPartUploadComplete(req.Key, req.UploadId, object)
	if err != nil {
		log.Printf("分片传输验证失败：%v", err)
		return nil, err
	}

	//把数据更新到存储中心表repository_pool中
	rp := new(models.RepositoryPool)
	rp.Identity = tools.GetUUID()
	rp.Hash = req.Md5
	rp.Name = req.Name
	rp.Size = req.Size
	rp.Ext = req.Ext
	rp.Path = define.TenCosURL + "/" + req.Key
	_, err = l.svcCtx.MysqlEngine.InsertOne(rp)
	if err != nil {
		log.Printf("分片上传文件插入存储中心表失败：%v", err)
		return nil, err
	}

	//同时更新用户文件关联表
	ur := new(models.UserRepository)
	ur.Identity = tools.GetUUID()
	ur.UserIdentity = userIdentity
	ur.RepositoryIdentity = rp.Identity
	ur.ParentId = req.ParentId
	ur.Ext = req.Ext
	ur.Name = req.Name
	_, err = l.svcCtx.MysqlEngine.InsertOne(ur)
	if err != nil {
		log.Printf("分片上传文件插入用户文件关联表失败：%v", err)
		return nil, err
	}

	//设置响应类型
	resp = new(types.FileUploadChunkCompleteResponse)
	resp.Identity = rp.Identity
	return resp, nil
}
