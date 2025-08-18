package connect

import (
	"context"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// Client 是一个可以复用，线程安全的HTTP客户端
	// 在整个应用程序中只初始化一次，然后复用
	Client *http.Client
)

// init 函数在这个包首次被调用的时候自动执行
func init() {
	// 创建一个默认的Transpost，但是可以自定义一些关键参数
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	Client = &http.Client{
		Transport: transport,
		// 这个Timeout是整个请求的超时时间，包括DNS解析、连接、重定向和传输数据
		Timeout: 5 * time.Second, //如果5秒内没有响应，就超时
	}
}

// Get 发送一个Get请求，判断是否成功访问(成功返回 200 OK)
// 接收一个 context.Context 以便上层调用着可以控制请求的取消或超时
func Get(ctx context.Context, url string) bool {
	// 使用 http.NewRequestWithContext 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// 如果URL格式本身有问题
		logx.Errorw("failed to create new http request",
			logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
		return false
	}

	// 使用 Client.Do(req) 来执行请求
	resp, err := Client.Do(req)
	if err != nil {
		// 网络层面的错误，如超时、DNS错误、连接被拒绝等。
		logx.Errorw("connect client.Get failed",
			logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
		return false
	}
	// 使用 defer 语句确保 resp.Body 一定会被关闭，否则资源一直无法释放
	defer resp.Body.Close()

	// 检查HTTP状态码
	return resp.StatusCode == http.StatusOK
}
