package logic

import (
	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"context"
	"log"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest) (resp *types.FileUploadResponse, err error) {
	//先判断该文件是否存在中心存储池中
	rp := new(models.RepositoryPool)
	_, err = l.svcCtx.MysqlEngine.Where("hash = ?", req.Hash).Get(rp)
	if err != nil {
		log.Printf("文件查询失败：%v", err)
		return nil, err
	}

	//从handler中拿出预处理过的数据（放在req中），然后把数据保存到RepositoryPool表中
	rp = &models.RepositoryPool{
		Identity: tools.GetUUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	//保存数据
	_, err = l.svcCtx.MysqlEngine.InsertOne(rp)
	if err != nil {
		log.Printf("用户上传数据记录RepositoryPool表失败：%v", err)
		return nil, err
	}

	//响应数据
	resp = &types.FileUploadResponse{
		Identity: rp.Identity,
		Ext:      rp.Ext,
		Name:     rp.Name,
	}

	return resp, nil
}
