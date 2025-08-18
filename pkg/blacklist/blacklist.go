package blacklist

var blackListSet map[string]struct{}

// SetBlackListSet 设置黑名单
func SetBlackListSet(list []string) {
	blackListSet = make(map[string]struct{})
	for _, v := range list {
		blackListSet[v] = struct{}{}
	}
}

// GetBlackListSet 获取黑名单
func GetBlackListSet()map[string]struct{}{
	return blackListSet 
}
