package handler

import (
	"net/http"

	"lanshan11/internal/logic"
	"lanshan11/internal/svc"
	"lanshan11/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ✅ 首字母大写，才能被 routes.go 引用
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
