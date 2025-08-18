package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shorturl/internal/svc"
	"shorturl/internal/types"
	"shorturl/model"
	"shorturl/pkg/base62"
	"shorturl/pkg/blacklist"
	"shorturl/pkg/connect"
	"shorturl/pkg/md5"
	"shorturl/pkg/urltools"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var shortUrl string

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链业务逻辑：输入一个长链接->转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 参数校验
	// 1.1. 参数规则校验 在handler/converandler.go里面已经实现
	// 1.2. 输入的长链接必须是一个能够请求通的网址
	if ok := connect.Get(l.ctx, req.LongUrl); !ok {
		return nil, errors.New("无效的链接")
	}
	// 1.3. 判断之前是否转链过（数据库中是否已存在长链接）
	// 将LongURL转为MD5字符串，然后去数据库中查看
	longUrlMD5 := md5.NewMD5String(req.LongUrl)
	u, err := l.svcCtx.ShortUrlMap.FindOneByMd5(l.ctx, sql.NullString{String: longUrlMD5, Valid: true})
	if err != sqlx.ErrNotFound { // 如果不是属于查询不到
		if err == nil { // 查询到了
			return nil, fmt.Errorf("该链接已经转为%s", u.Surl.String)
		}

		logx.Errorw("ShortUrlMap.FindOneByMd5 failed", // 其他类型的错误
			logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
		return nil, err
	}
	// 到达这里说明并没有被转链过
	// 1.4. 输入的不能是一个短链接（避免循环转链）
	base, err := urltools.GetURLPathBase(req.LongUrl)
	if err != nil {
		logx.Errorw("GetURLPathBase failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, err
	}
	if _, err := l.svcCtx.ShortUrlMap.FindOneBySurl(l.ctx, sql.NullString{String: base, Valid: true}); err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errors.New("该链接已经是短链接了")
		}
		logx.Errorw("ShortUrlMap.FindOneBySurl failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, err
	}

	for {
		// 2. 取号
		// 到了这里，说明已经这个长链接既没有被转链过，也不是一个短链接
		res, err := l.svcCtx.Sequence.ReplaceIntoByStub(l.ctx, "a")
		if err != nil {
			logx.Errorw("Sequence.ReplaceIntoByStub failed", logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
			return nil, err
		}
		// 获取号码
		num, err := res.LastInsertId()
		if err != nil {
			logx.Errorw("res.LastInsertId failed", logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
			return nil, err
		}
		// 3. 号码转短链
		shortUrl = base62.ToBase62(uint64(num))

		// 如果不是在黑名单的才退出
		if _, ok := blacklist.GetBlackListSet()[shortUrl]; !ok {
			break
		}
	}

	// 4. 存储长链接映射关系
	_, err = l.svcCtx.ShortUrlMap.Insert(l.ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongUrl, Valid: true},
		Md5:  sql.NullString{String: longUrlMD5, Valid: true},
		Surl: sql.NullString{String: shortUrl, Valid: true},
	})
	if err != nil {
		logx.Errorw("ShortUrlMap.Insert failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, err
	}
	// 5. 返回短链接
	return &types.ConvertResponse{
		ShortUrl: shortUrl,
	}, nil
}
