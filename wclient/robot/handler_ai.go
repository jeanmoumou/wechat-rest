package robot

import (
	"strings"

	"github.com/opentdp/wechat-rest/dbase/profile"
	"github.com/opentdp/wechat-rest/wcferry"
	"github.com/opentdp/wechat-rest/wclient/aichat"
)

func aiHandler() {

	handlers["/ai"] = &Handler{
		Level:    0,
		ChatAble: true,
		RoomAble: true,
		Describe: "提问或交谈",
		Callback: func(msg *wcferry.WxMsg) string {
			text := strings.TrimSpace(msg.Content)
			if text != "" {
				return aichat.Text(msg.Sender, msg.Roomid, text)
			}
			return "请在指令后输入问题"
		},
	}

}

func aiMessagePrefix(msg *wcferry.WxMsg) string {

	if len(msg.Content) == 0 {
		return ""
	}

	if msg.Content[0:1] != "/" {
		if strings.Contains(msg.Xml, self().Wxid) {
			msg.Content = "/ai " + msg.Content
		} else {
			up, _ := profile.Fetch(&profile.FetchParam{Wxid: msg.Sender, Roomid: msg.Roomid})
			if up.AiArgot == "" {
				if !msg.IsGroup {
					msg.Content = "/ai " + msg.Content
				}
			} else if strings.HasPrefix(msg.Content, up.AiArgot) {
				msg.Content = strings.Replace(msg.Content, up.AiArgot, "/ai ", 1)
			}
		}
	}

	return ""

}
