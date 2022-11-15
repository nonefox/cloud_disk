package handler

import (
	"net/http"

	"cloud_disk/greet/internal/logic"
	"cloud_disk/greet/internal/svc"
	"cloud_disk/greet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserRepositoryRelationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRepositoryRelationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUserRepositoryRelationLogic(r.Context(), svcCtx)
		//在auth_middleware用户验证中，在*http.Request中我们存入了解析token中的userClaim中的用户信息
		resp, err := l.UserRepositoryRelation(&req, r.Header.Get("UserIdentity"))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
