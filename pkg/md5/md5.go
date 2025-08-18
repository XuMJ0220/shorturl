package md5

import (
	"crypto/md5"
	"fmt"
)

func NewMD5String(str string) string{
	data:=[]byte(str)
	return fmt.Sprintf("%x",md5.Sum(data))
}