package svc

import (
	"shorturl/internal/config"
	"shorturl/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	ShortUrlMap model.ShortUrlMapModel
	Sequence    model.SequenceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn1 := sqlx.NewMysql(c.ShortUrlMapMysql.DSN)
	conn2 := sqlx.NewMysql(c.Sequence.DSN)
	return &ServiceContext{
		ShortUrlMap: model.NewShortUrlMapModel(conn1),
		Sequence:    model.NewSequenceModel(conn2),
		Config:      c,
	}
}
