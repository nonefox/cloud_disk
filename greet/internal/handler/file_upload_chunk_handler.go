package handler

import (
	"cloud_disk/greet/tools"
	"errors"
	"net/http"

	"cloud_disk/greet/internal/logic"
	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadChunkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		//获取从前端传过来的formdata的数据
		if r.PostForm.Get("key") == "" {
			httpx.Error(w, errors.New("key不能为空"))
			return
		}
		if r.PostForm.Get("upload_id") == "" {
			httpx.Error(w, errors.New("upload_id不能为空"))
			return
		}
		if r.PostForm.Get("part_number") == "" {
			httpx.Error(w, errors.New("part_number不能为空"))
			return
		}

		//获取文件MD5值
		etag, err := tools.TenCosPartUpload(r)
		if err != nil {

			return
		}
		l := logic.NewFileUploadChunkLogic(r.Context(), svcCtx)
		resp, err := l.FileUploadChunk(&req)
		//放入响应结构中
		resp = new(types.FileUploadChunkResponse)
		resp.Etag = etag
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
