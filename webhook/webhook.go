package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jummyliu/pkg/request"
)

const (
	WEBHOOK_TYPE_WEIXIN_WORK string = "企微机器人"
	WEBHOOK_TYPE_DINGDING    string = "钉钉机器人"
	WEBHOOK_TYPE_FEISHU      string = "飞书机器人"
	MESSAGE_TYPE_TEXT        string = "text"
	MESSAGE_TYPE_MARKDOWN    string = "markdown"
	MESSAGE_TYPE_IMAGE       string = "image"
)

var DICT_WEBHOOK_TYPE_TO_SERVER_ADDRESS = map[string]string{
	WEBHOOK_TYPE_WEIXIN_WORK: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=",
	WEBHOOK_TYPE_DINGDING:    "https://oapi.dingtalk.com/robot/send?access_token=",
	WEBHOOK_TYPE_FEISHU:      "https://open.feishu.cn/open-apis/bot/v2/hook/",
}

var DICT_WEBHOOK_TYPE_TO_MESSAGE_TEMPLATE = map[string]map[string]map[string]any{
	WEBHOOK_TYPE_WEIXIN_WORK: {
		MESSAGE_TYPE_TEXT:     {"msgtype": "text", "text": map[string]string{"content": ""}},
		MESSAGE_TYPE_MARKDOWN: {"msgtype": "markdown", "markdown": map[string]string{"content": ""}},
		MESSAGE_TYPE_IMAGE:    {"msgtype": "image", "image": map[string]string{"base64": "", "md5": ""}},
	},
	WEBHOOK_TYPE_DINGDING: {
		MESSAGE_TYPE_TEXT:     {"msgtype": "text", "text": map[string]string{"content": ""}},
		MESSAGE_TYPE_MARKDOWN: {"msgtype": "markdown", "markdown": map[string]string{"title": "新消息", "text": ""}},
		// MESSAGE_TYPE_IMAGE:    {"msgtype": "markdown", "markdown": map[string]string{"title": "新消息", "text": "![]({{image_url}})"}}, // 钉钉机器人不支持 image 类型，改为使用 markdown 类型
	},
	WEBHOOK_TYPE_FEISHU: {
		MESSAGE_TYPE_TEXT:     {"msg_type": "text", "content": map[string]string{"text": ""}},
		MESSAGE_TYPE_MARKDOWN: {"msg_type": "text", "content": map[string]string{"text": ""}}, // 飞书机器人不支持 markdown 类型，改为使用 text 类型
		// MESSAGE_TYPE_IMAGE:    {"msg_type": "image", "content": map[string]string{"image_key": ""}}, // 飞书机器人发送图片需要先上传至 飞书开放平台 获取 image_key
	},
}

type WebhookSender struct {
	api_key         string
	server_address  string
	webhook_type    string
	proxy_address   string
	message_title   string
	request_headers map[string]string
}

type OptionFunc func(*WebhookSender)

func initOptions(options ...OptionFunc) *WebhookSender {
	webhook_sender := &WebhookSender{
		api_key:        "",
		server_address: "",
		webhook_type:   WEBHOOK_TYPE_WEIXIN_WORK,
		proxy_address:  "",
		message_title:  "",
		request_headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
			"Accept":          "application/json",
			"Content-Type":    "application/json;charset=UTF-8",
			"Accept-Encoding": "gzip, deflate",
			"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
			"Connection":      "close",
		},
	}
	for _, option_func := range options {
		option_func(webhook_sender)
	}
	return webhook_sender
}

func WithOptions(options WebhookSender) OptionFunc {
	return func(webhook_sender *WebhookSender) {
		*webhook_sender = options
	}
}

// 可指定服务端地址，为空时使用 webhook 类型的默认地址
func WithServerAddress(s string) OptionFunc {
	return func(webhook_sender *WebhookSender) {
		webhook_sender.server_address = s
	}
}

// 可指定 webhook 类型，为空时默认使用 企微机器人
func WithWebhookType(s string) OptionFunc {
	return func(webhook_sender *WebhookSender) {
		webhook_sender.webhook_type = s
	}
}

// 可指定代理地址，为空时不使用代理
func WithProxyAddress(s string) OptionFunc {
	return func(webhook_sender *WebhookSender) {
		webhook_sender.proxy_address = s
	}
}

// 可指定请求头，为空时使用默认请求头
func WithRequestHeaders(m map[string]string) OptionFunc {
	return func(webhook_sender *WebhookSender) {
		webhook_sender.request_headers = m
	}
}

// 钉钉机器人的 Markdown 类型支持指定消息标题，不指定时默认使用“新消息”作为标题
func WithMessageTitle(s string) OptionFunc {
	return func(webhook_sender *WebhookSender) {
		webhook_sender.message_title = s
	}
}

func New(api_key string, options ...OptionFunc) *WebhookSender {
	webhook_sender := initOptions(options...)
	webhook_sender.api_key = api_key
	// 服务端地址为空时，使用 webhook 类型的默认地址
	if webhook_sender.server_address == "" {
		server_address, ok := DICT_WEBHOOK_TYPE_TO_SERVER_ADDRESS[webhook_sender.webhook_type]
		if ok {
			webhook_sender.server_address = server_address
		} else {
			webhook_sender.server_address = DICT_WEBHOOK_TYPE_TO_SERVER_ADDRESS[WEBHOOK_TYPE_WEIXIN_WORK]
		}
	}

	return webhook_sender
}

func (webhook_sender *WebhookSender) SendMessage(content map[string]any) (err error) {
	request_data, _ := json.Marshal(content)
	_, _, _, err = request.DoRequest(
		fmt.Sprintf("%s%s", webhook_sender.server_address, webhook_sender.api_key),
		request.WithMethod(http.MethodPost),
		request.WithHeader(webhook_sender.request_headers),
		request.WithData(request_data),
		request.WithProxy(webhook_sender.proxy_address),
	)
	if err != nil {
		return
	}
	return
}

// 推送Text类型消息
func (webhook_sender *WebhookSender) SendMessageText(content string) (err error) {
	map_content := DICT_WEBHOOK_TYPE_TO_MESSAGE_TEMPLATE[webhook_sender.webhook_type][MESSAGE_TYPE_TEXT]
	var temp_content = map[string]string{}
	switch webhook_sender.webhook_type {
	case WEBHOOK_TYPE_WEIXIN_WORK:
	case WEBHOOK_TYPE_DINGDING:
		temp_content = map_content["text"].(map[string]string)
		temp_content["content"] = content
		map_content["text"] = temp_content
	case WEBHOOK_TYPE_FEISHU:
		temp_content = map_content["content"].(map[string]string)
		temp_content["text"] = content
		map_content["content"] = temp_content
	}
	err = webhook_sender.SendMessage(map_content)
	return
}

// 推送Markdown类型消息
//
//	飞书机器人不支持Markdown格式降级为text类型
//	钉钉机器人的 Markdown 类型支持指定消息标题，不指定时默认使用“新消息”作为标题
func (webhook_sender *WebhookSender) SendMessageMarkdown(content string) (err error) {
	map_content := DICT_WEBHOOK_TYPE_TO_MESSAGE_TEMPLATE[webhook_sender.webhook_type][MESSAGE_TYPE_MARKDOWN]
	var temp_content = map[string]string{}
	switch webhook_sender.webhook_type {
	case WEBHOOK_TYPE_WEIXIN_WORK:
		temp_content = map_content["markdown"].(map[string]string)
		temp_content["content"] = content
		map_content["markdown"] = temp_content
	case WEBHOOK_TYPE_DINGDING:
		temp_content = map_content["markdown"].(map[string]string)
		temp_content["text"] = content
		if webhook_sender.message_title != "" {
			temp_content["title"] = webhook_sender.message_title
		}
		map_content["markdown"] = temp_content
	case WEBHOOK_TYPE_FEISHU:
		temp_content = map_content["content"].(map[string]string)
		temp_content["text"] = content
		map_content["content"] = temp_content
	}
	err = webhook_sender.SendMessage(map_content)
	return
}

func SendToWeiXinWork(api_key string, content map[string]any, proxy_address string) (err error) {

	request_url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + api_key

	request_data, _ := json.Marshal(content)

	request_headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
		"Accept":          "application/json",
		"Content-Type":    "application/json;charset=UTF-8",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Connection":      "close",
	}

	_, _, _, err = request.DoRequest(
		request_url,
		request.WithMethod(http.MethodPost),
		request.WithHeader(request_headers),
		request.WithData(request_data),
		request.WithProxy(proxy_address),
	)

	if err != nil {
		return
	}

	return
}

func SendToDingDing(api_key string, content map[string]any, proxy_address string) (err error) {

	request_url := "https://oapi.dingtalk.com/robot/send?access_token=" + api_key

	request_data, _ := json.Marshal(content)

	request_headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
		"Accept":          "application/json, text/plain, */*",
		"Content-Type":    "application/json;charset=UTF-8",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Connection":      "close",
	}

	_, _, _, err = request.DoRequest(
		request_url,
		request.WithMethod(http.MethodPost),
		request.WithHeader(request_headers),
		request.WithData(request_data),
		request.WithProxy(proxy_address),
	)

	if err != nil {
		return
	}

	return
}

func SendToFeiShu(api_key string, content map[string]any, proxy_address string) (err error) {

	request_url := "https://open.feishu.cn/open-apis/bot/v2/hook/" + api_key

	request_data, _ := json.Marshal(content)

	request_headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
		"Accept":          "application/json, text/plain, */*",
		"Content-Type":    "application/json;charset=UTF-8",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Connection":      "close",
	}

	_, _, _, err = request.DoRequest(
		request_url,
		request.WithMethod(http.MethodPost),
		request.WithHeader(request_headers),
		request.WithData(request_data),
		request.WithProxy(proxy_address),
	)

	if err != nil {
		return
	}

	return
}
