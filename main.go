package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/2mf8/Bot-Client-Go/safe_ws"

	bot "github.com/2mf8/Better-Bot-Go"
	bytesimage "github.com/2mf8/Better-Bot-Go/bytes_image"
	"github.com/2mf8/Better-Bot-Go/dto"
	"github.com/2mf8/Better-Bot-Go/dto/keyboard"
	"github.com/2mf8/Better-Bot-Go/openapi"
	"github.com/2mf8/Better-Bot-Go/token"
	"github.com/2mf8/Better-Bot-Go/webhook"
	log "github.com/sirupsen/logrus"
)

var Apis = make(map[string]openapi.OpenAPI, 0)

func main() {
	safe_ws.InitLog()
	as := webhook.ReadSetting()
	for i, v := range as.Apps {
		go safe_ws.ConnectUniversal(fmt.Sprintf("%v", v.AppId), v.WSSAddr)
		token := token.BotToken(v.AppId, v.Token, string(token.TypeBot))
		api := bot.NewOpenAPI(token).WithTimeout(3 * time.Second)
		Apis[i] = api
	}
	b, _ := json.Marshal(as)
	fmt.Println("配置", string(b))
	safe_ws.GroupAtMessageEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		ctx := context.WithValue(context.Background(), "key", "value")
		content := strings.TrimSpace(data.Content)
		log.Info(data.Content, data.GroupId, " <- ", content)
		if content == "base" {
			s, err := bytesimage.GetImageBytes("http://2mf8.cn:2014/view/333.png?scramble=R")
			fmt.Println(string(s))
			if err == nil {
				resp, err := Apis[appid].PostGroupRichMediaMessage(ctx, data.GroupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
				fmt.Println(err)
				if resp != nil {
					newMsg := &dto.GroupMessageToCreate{
						Media: &dto.FileInfo{
							FileInfo: resp.FileInfo,
						},
						MsgID:   data.MsgId,
						MsgType: 7,
						MsgReq:  1,
					}
					Apis[appid].PostGroupMessage(ctx, data.GroupId, newMsg)
				}
			}
		}
		if content == "kb" {
			ctx := context.WithValue(context.Background(), "key", "value")
			/* rows := keyboard.CustomKeyboard{} */
			/* kb := gkb.Builder().
			TextButton("测试", "已测试", "成功", false, true).
			UrlButton("爱魔方吧", "一仝", "https://2mf8.cn", false, true).
			SetRow().
			TextButton("测试", "已测试", "成功", false, true).
			SetRow()
			b, _:= json.Marshal(kb)
			json.Unmarshal(b, &rows) */
			fmt.Println("测试")
			Apis[appid].PostGroupMessage(ctx, data.GroupId, &dto.C2CMessageToCreate{
				Keyboard: &keyboard.MessageKeyboard{
					ID: "101981675_1735044770",
				},
				MsgType: dto.C2CMsgTypeMarkdown,
				MsgID:   data.MsgId,
			})
		}
		return nil
	}
	safe_ws.C2CMessageEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		content := strings.TrimSpace(data.Content)
		log.Info(data.Content, data.Author.UserOpenId, " -", content)
		if content == "kb" {
			ctx := context.WithValue(context.Background(), "key", "value")
			/* rows := keyboard.CustomKeyboard{} */
			/* kb := gkb.Builder().
			TextButton("测试", "已测试", "成功", false, true).
			UrlButton("爱魔方吧", "一仝", "https://2mf8.cn", false, true).
			SetRow().
			TextButton("测试", "已测试", "成功", false, true).
			SetRow()
			b, _:= json.Marshal(kb)
			json.Unmarshal(b, &rows) */
			fmt.Println("测试")
			s, err := bytesimage.GetImageBytes("http://2mf8.cn:2014/view/333.png?scramble=R")
			resp, err := Apis[appid].PostC2CRichMediaMessage(ctx, data.Author.UserOpenId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
			fmt.Println(err)
			if resp != nil {
				newMsg := &dto.C2CMessageToCreate{
					Media: &dto.FileInfo{
						FileInfo: resp.FileInfo,
					},
					MsgID:   data.Id,
					MsgType: dto.C2CMsgTypeMedia,
					MsgReq:  1,
				}
				Apis[appid].PostC2CMessage(ctx, data.Author.UserOpenId, newMsg)
			}
			Apis[appid].PostC2CMessage(ctx, data.Author.UserOpenId, &dto.C2CMessageToCreate{
				Keyboard: &keyboard.MessageKeyboard{
					ID: "101981675_1735044770",
				},
				MsgType: dto.C2CMsgTypeMarkdown,
				MsgID:   data.Id,
			})
		}
		return nil
	}
	safe_ws.MessageEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSMessageData) error {
		s, _ := bytesimage.GetImageBytes("./333.png")
		_, e := Apis[appid].PostFormFileReaderImage(context.Background(), data.ChannelID, map[string]string{
			"msg_id":  data.ID,
			"content": "333.png",
		}, "333.png", bytes.NewBuffer(s))
		fmt.Println(e)
		return nil
	}
	safe_ws.InteractionEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSInteractionData) error {
		/* fmt.Println(data.ChannelID)
		ctx := context.WithValue(context.Background(), "key", "value")
			Apis[appid].PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{
				Content: "测试",
				MsgID: data.ID,
			}) */
		return nil
	}
	safe_ws.FriendAddEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSFriendAddData) error {
		Apis[appid].PostC2CMessage(context.Background(), data.OpenId, &dto.C2CMessageToCreate{
			Content: "hello",
			EventID: dto.EventType(event.ID),
		})
		return nil
	}
	safe_ws.GroupAddRobotEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSGroupAddRobotData) error {
		fmt.Println(data.GroupOpenId, data.OpMemberOpenId, data.Timestamp)
		m, e := Apis[appid].PostGroupMessage(context.Background(), data.GroupOpenId, &dto.C2CMessageToCreate{
			Content: "hello",
			EventID: dto.EventType(event.ID),
		})
		fmt.Println(m,e)
		return nil
	}
	select {}
}
