package urltools

import (
	"net/url"
	"path"
)

func GetURLPathBase(targetUrl string) (string, error) {
	// 将传进来的rul进行解析
	u, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	}
	// 获取Path部分,并取得Path的最后一个斜杠后的部分
	base := path.Base(u.Path)
	return base, nil
}
