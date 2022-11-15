package middleware

import (
	"cloud_disk/greet/tools"
	"net/http"
	"strconv"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// Handle 用户验证
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//获取用户的认证信息(也就是我们的token信息)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		//我们解析用户token，获取出里面的userClaim信息
		userClaim, err := tools.AnalyseToken(auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//解析之后把用户声明信息，放回*http.Request的header中(方便后面取用)
		r.Header.Set("UserId", strconv.Itoa(userClaim.Id))
		r.Header.Set("UserIdentity", userClaim.Identity)
		r.Header.Set("UserName", userClaim.Name)

		next(w, r)
	}
}
