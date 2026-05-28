package handler

import (
	"net/http"

	"lanshan11/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/register",
				Handler: RegisterHandler(serverCtx), // ✅ 现在能找到了
			},
			// 后面可以加 login/refresh 等路由
		},
	)
}
