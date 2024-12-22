package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/2mf8/Bot-Client-Go/safe_ws"

	bot "github.com/2mf8/Better-Bot-Go"
	"github.com/2mf8/Better-Bot-Go/dto"
	"github.com/2mf8/Better-Bot-Go/openapi"
	"github.com/2mf8/Better-Bot-Go/token"
	"github.com/2mf8/Better-Bot-Go/webhook"
	log "github.com/sirupsen/logrus"
)

var Apis = make(map[string]openapi.OpenAPI, 0)

func main() {
	go safe_ws.ConnectUniversal()
	safe_ws.InitLog()
	as := webhook.ReadSetting()
	for i, v := range as.Apps {
		token := token.BotToken(v.AppId, v.Token, string(token.TypeBot))
		api := bot.NewSandboxOpenAPI(token).WithTimeout(3 * time.Second)
		Apis[i] = api
	}
	b, _ := json.Marshal(as)
	fmt.Println("配置", string(b))
	safe_ws.GroupAtMessageEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		log.Info(data.Content, data.GroupId)
		ctx := context.WithValue(context.Background(), "key", "value")
		newMsg := &dto.GroupMessageToCreate{
			Content: "测试",
			MsgID:   data.MsgId,
			MsgType: 0,
		}
		Apis[appid].PostGroupMessage(ctx, data.GroupId, newMsg)
		return nil
	}
	select {}
}
