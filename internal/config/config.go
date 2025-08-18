package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	ShortUrlMapMysql struct {
		DSN string
	}

	Sequence struct {
		DSN string
	}

	Base62CharacterSet string

	BlackList []string
}
