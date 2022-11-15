package handler

import (
	"cloud_disk/greet/internal/logic"
	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"
	"cloud_disk/greet/models"
	"cloud_disk/greet/tools"
	"crypto/md5"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"log"
	"net/http"
	"path"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		//在这里对前端传过来的数据进行预处理
		file, header, err := r.FormFile("file") //获取前端传过来的文件
		if err != nil {
			httpx.Error(w, err)
			return
		}
		//与用户存储中心的数据进行比对，看用户是否已经上传该文件
		fileBytes := make([]byte, header.Size)
		_, err = file.Read(fileBytes)
		if err != nil {
			log.Printf("读取数据失败：%v", err)
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(fileBytes))
		//用这个值来验证是否文件是否存在RepositoryPool表中
		rp := new(models.RepositoryPool)
		get, err := svcCtx.MysqlEngine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			return
		}
		if get {
			//如果存在，返回该文件的Identity,Ext,Name
			httpx.OkJson(w, &types.FileUploadResponse{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
		}

		//如果不存在，我们就把他上传到COS
		fileURL, err := tools.TenCosUpload(r)
		if err != nil {
			return
		}

		//然后把上传文件的信息放到我们的req中(传到logic中，让他去保存上传数据信息)
		req.Name = header.Filename
		req.Size = header.Size
		req.Ext = path.Ext(header.Filename) //文件后缀名
		req.Hash = hash
		req.Path = fileURL //文件访问路径

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
