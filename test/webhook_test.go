package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/webhook"
)

func TestNewWebhookSender(t *testing.T) {
	api_key := "your_api_key"

	// 样例一：初始化对象，可重复调用函数发送消息
	webhook_sender := webhook.New(
		api_key,
		webhook.WithWebhookType(webhook.WEBHOOK_TYPE_WEIXIN_WORK),
		// webhook.WithWebhookType(webhook.WEBHOOK_TYPE_DINGDING),
		// webhook.WithWebhookType(webhook.WEBHOOK_TYPE_FEISHU),
	)

	webhook_sender.SendMessageText("test1")
	webhook_sender.SendMessageText("test2")

	// 样例二：实例化对象后直接调用函数发送消息
	content := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": "test—",
		},
	}
	webhook.New(
		api_key,
		webhook.WithWebhookType(webhook.WEBHOOK_TYPE_WEIXIN_WORK),
	).SendMessage(content)

}

func TestSendToWeiXinWork(t *testing.T) {
	api_key := "your_api_key"
	content := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": "test",
		},
	}
	err := webhook.SendToWeiXinWork(api_key, content, "")
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}
}

// 钉钉机器人推送Markdown消息必须指定标题
func TestSendToDingDing(t *testing.T) {
	api_key := "your_api_key"
	content := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "test",
			"text":  "test",
		},
	}

	err := webhook.SendToDingDing(api_key, content, "")
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}
}

func TestSendToFeiShu(t *testing.T) {
	api_key := "your_api_key"
	content := map[string]any{
		"msg_type": "text",
		"content": map[string]string{
			"text": "test",
		},
	}
	err := webhook.SendToFeiShu(api_key, content, "")
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}
}
