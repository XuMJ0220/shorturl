package logic

import (
	"context"
	"database/sql"

	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 根据端连接，从数据库中查询
	u, err := l.svcCtx.ShortUrlMap.FindOneBySurl(l.ctx, sql.NullString{
		String: req.ShortUrl,
		Valid:  true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		logx.Errorw("ShortUrlMap.FindOneBySurl failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, err
	}
	longUrl := u.Lurl.String
	// 2. 返回长链接，进行重映射
	return &types.ShowResponse{
		LongUrl: longUrl,
	}, nil
}
