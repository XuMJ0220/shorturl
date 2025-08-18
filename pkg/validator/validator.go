package validator

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
)

// 定义一个包级别的私有变量来持有单例
var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator 获取单例
func GetValidator() *validator.Validate {
	// 使用 sync.Once 确保在多协程环境只执行一次初始化
	once.Do(func() {
		// 在这里执行初始化操作
		validate = validator.New()

		// Do 下面执行其他自己想要的操作
	})

	return validate
}

// Struct 是一个方便的函数，直接使用单例验证器来验证结构体
func Struct(s any) error {
	return GetValidator().Struct(s)
}

// StructCtx 是一个方便的函数，直接使用单例验证器来验证结构体,可以传入上下文
func StructCtx(ctx context.Context, s any) error {
	return GetValidator().StructCtx(ctx, s)
}

// Var 是一个方便的函数，用于验证单个变量
func Var(field interface{}, tag string) error {
	return GetValidator().Var(field, tag)
}
