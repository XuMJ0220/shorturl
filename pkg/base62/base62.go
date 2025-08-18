package base62

import "strings"

var base62CharacterSet string 

var baseLength uint64

// ToBase62 将一个十进制数转换为base62字符串
func ToBase62(n uint64) string {
	if n == 0 {
		return "0"
	}

	var result strings.Builder
	result.Grow(11)

	for n > 0 {
		remainder := n % baseLength
		result.WriteByte(base62CharacterSet[remainder])
		n /= baseLength
	}

	//因为我们也需要考虑到转换的安全，所以这里没必要对result.String()进行倒转
	return result.String()
}

// SetCharacterSet 设置base62字符集
func SetCharacterSet(charset string) {
	base62CharacterSet = charset
	baseLength = uint64(len(base62CharacterSet))
}
