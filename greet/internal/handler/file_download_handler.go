package handler

import (
	"net/http"

	"cloud_disk/greet/internal/logic"
	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDownloadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewFileDownloadLogic(r.Context(), svcCtx)
		resp, err := l.FileDownload(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
