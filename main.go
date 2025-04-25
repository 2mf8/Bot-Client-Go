package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	bot "github.com/2mf8/Better-Bot-Go"
	bytesimage "github.com/2mf8/Better-Bot-Go/bytes_image"
	"github.com/2mf8/Better-Bot-Go/dto"
	v1 "github.com/2mf8/Better-Bot-Go/openapi/v1"
	"github.com/2mf8/Better-Bot-Go/token"
	"github.com/2mf8/Better-Bot-Go/webhook"
	"github.com/2mf8/Bot-Client-Go/safe_ws"
	log "github.com/sirupsen/logrus"
)

func main() {
	safe_ws.InitLog()
	as := webhook.ReadSetting()
	for _, v := range as.Apps {
		atr := v1.GetAccessToken(fmt.Sprintf("%v", v.AppId), v.AppSecret)
		iat, err := strconv.Atoi(atr.ExpiresIn)
		if err == nil && atr.AccessToken != "" {
			aei := time.Now().Unix() + int64(iat)
			token := token.BotToken(v.AppId, atr.AccessToken, string(token.TypeQQBot))
			if v.IsSandBox {
				api := bot.NewSandboxOpenAPI(token).WithTimeout(3 * time.Second)
				go bot.AuthAcessAdd(fmt.Sprintf("%v", v.AppId), &bot.AccessToken{AccessToken: atr.AccessToken, ExpiresIn: aei, Api: api, AppSecret: v.AppSecret, IsSandBox: v.IsSandBox, Appid: v.AppId})
			} else {
				api := bot.NewOpenAPI(token).WithTimeout(3 * time.Second)
				go bot.AuthAcessAdd(fmt.Sprintf("%v", v.AppId), &bot.AccessToken{AccessToken: atr.AccessToken, ExpiresIn: aei, Api: api, AppSecret: v.AppSecret, IsSandBox: v.IsSandBox, Appid: v.AppId})
			}
		}
		time.Sleep(time.Millisecond * 100)
		if as.IsOpen {
			go safe_ws.ConnectUniversalWithSecret(fmt.Sprintf("%v", v.AppId), v.AppSecret, v.WSSAddr)
		} else {
			go safe_ws.ConnectUniversal(fmt.Sprintf("%v", v.AppId), v.WSSAddr)
		}
	}
	/* b, _ := json.Marshal(as)
	fmt.Println("配置", string(b)) */
	safe_ws.GroupAtMessageEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		ctx := context.WithValue(context.Background(), "key", "value")
		content := strings.TrimSpace(data.Content)
		log.Info(data.Content, data.GroupId, " <- ", content)
		if content == "测试" {
			newMsg := &dto.GroupMessageToCreate{
				Content: "测试多处WSS发消息成功",
				MsgType: dto.C2CMsgTypeText,
				MsgID:   data.MsgId,
				MsgReq:  4,
			}
			g, e := bot.SendApi(appid).PostGroupMessage(ctx, data.GroupId, newMsg)
			fmt.Println(g.Id, e)
		}
		if content == "base" {
			s, err := bytesimage.GetImageBytes("http://2mf8.cn:2014/view/333.png?scramble=R")
			fmt.Println(string(s))
			if err == nil {
				resp, err := bot.SendApi(appid).PostGroupRichMediaMessage(ctx, data.GroupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, FileData: s, SrvSendMsg: false})
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
					bot.SendApi(appid).PostGroupMessage(ctx, data.GroupId, newMsg)
				}
			}
		}
		return nil
	}
	safe_ws.FriendAddEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSFriendAddData) error {
		bot.SendApi(appid).PostC2CMessage(context.Background(), data.OpenId, &dto.C2CMessageToCreate{
			Content: "hello",
			EventID: dto.EventType(event.ID),
		})
		return nil
	}
	safe_ws.GroupAddRobotEventHandler = func(appid string, event *dto.WSPayload, data *dto.WSGroupAddRobotData) error {
		fmt.Println(data.GroupOpenId, data.OpMemberOpenId, data.Timestamp)
		m, e := bot.SendApi(appid).PostGroupMessage(context.Background(), data.GroupOpenId, &dto.C2CMessageToCreate{
			Content: "hello",
			EventID: dto.EventType(event.ID),
		})
		fmt.Println(m, e)
		return nil
	}
	select {}
}
