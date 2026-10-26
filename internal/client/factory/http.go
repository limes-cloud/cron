package factory

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	json "github.com/json-iterator/go"
	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/client/service"
)

type HttpRequest struct {
	URL      string            `json:"url"`
	Params   map[string]any    `json:"params"`
	Method   string            `json:"method"`
	Header   map[string]string `json:"header"`
	Body     string            `json:"body"`
	BodyJson map[string]any    `json:"bodyJson"`
}

func (f *Factory) http(ctx kratosx.Context, task *service.ExecTaskRequest) (int, error) {
	var (
		data HttpRequest
		code = defaultErrorCode
		body string
		err  error
	)
	if err = json.Unmarshal([]byte(task.Value), &data); err != nil {
		return code, err
	}
	if data.Body != "" {
		body = data.Body
	}
	if len(data.BodyJson) != 0 {
		body, err = json.MarshalToString(data.BodyJson)
		if err != nil {
			return code, err
		}
	}

	// 使用 resty 直接构建请求：kratosx 封装的 response 未暴露 StatusCode，无法与 ExpectCode 匹配。
	// 重试语义由 server 端 RetryCount 循环统一控制，此处不再重复设置 resty retry，避免执行次数平方。
	client := resty.New()
	if task.MaxExecTime != 0 {
		client.SetTimeout(time.Duration(task.MaxExecTime) * time.Second)
	} else {
		client.SetTimeout(60 * time.Second)
	}

	req := client.R()
	req.Header.Set("User-Agent", "cron-http-executor")
	for key, val := range data.Header {
		req.Header.Set(key, val)
	}
	for key, val := range data.Params {
		req.QueryParam.Set(key, fmt.Sprint(val))
	}
	req.SetBody(body)
	req.Method = data.Method
	req.URL = data.URL

	response, err := req.Send()
	if err != nil {
		return defaultErrorCode, err
	}

	f.reply(task.Uuid, logInfo, string(response.Body()))
	return response.StatusCode(), nil
}