package handler

import (
	"net/http"

	"shorturl/internal/logic"
	"shorturl/internal/svc"
	"shorturl/internal/types"
	"shorturl/pkg/validator"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 参数规则校验
		// 如果参数校验没通过，直接返回，不继续往下执行
		if err := validator.GetValidator().StructCtx(r.Context(), &req); err != nil {
			logx.Error("validator check failed", logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
			// 这里要向客户端返回错误信息，否则默认返回200状态码
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
