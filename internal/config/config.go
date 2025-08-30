package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlMapMysql struct {
		DSN string
	}

	Sequence struct {
		DSN string
	}

	CacheRedis cache.CacheConf

	Base62CharacterSet string

	BlackList []string
}
