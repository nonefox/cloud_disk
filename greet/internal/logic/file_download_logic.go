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

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadRequest) (resp *types.FileDownloadResponse, err error) {
	ur := new(models.RepositoryPool)
	//通过用户传过来的文件的identity，来查询user_repository表，获取出文件的path
	get, err := l.svcCtx.MysqlEngine.Where("identity = ?", req.Identity).Get(ur)
	if err != nil {
		log.Printf("查询文件数据失败：%v", err)
		return nil, err
	}
	if !get {
		return nil, errors.New("文件不存在")
	}

	//用你自己的桶的前缀
	beforePath := "https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com/cloud_disk/"
	//我们把获取到的文件的path，截取出文件的key
	key := ur.Path[len(beforePath):]
	fileName := ur.Name + ur.Ext
	//开始下载文件
	fileDir, err := tools.TenCosDownload(key, fileName)
	if err != nil {
		return nil, errors.New("文件下载失败")
	}
	//返回用户下载文件的位置
	resp = new(types.FileDownloadResponse)
	resp.FileDir = fileDir
	return resp, nil
}
