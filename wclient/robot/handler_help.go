package robot

import (
	"fmt"
	"strings"

	"github.com/opentdp/wechat-rest/dbase/keyword"
	"github.com/opentdp/wechat-rest/dbase/llmodel"
	"github.com/opentdp/wechat-rest/dbase/profile"
	"github.com/opentdp/wechat-rest/dbase/setting"
	"github.com/opentdp/wechat-rest/wcferry"
	"github.com/opentdp/wechat-rest/wclient/aichat"
)

func helpHandler() []*Handler {

	cmds := []*Handler{}

	cmds = append(cmds, &Handler{
		Level:    0,
		Order:    900,
		ChatAble: true,
		RoomAble: true,
		Command:  "/help",
		Describe: "查看帮助信息",
		Callback: helpCallback,
	})

	return cmds

}

func helpCallback(msg *wcferry.WxMsg) string {

	up, _ := profile.Fetch(&profile.FetchParam{Wxid: msg.Sender, Roomid: prid(msg)})

	// 别名映射表
	aliasMap := map[string]map[string]string{}
	keywords, err := keyword.FetchAll(&keyword.FetchAllParam{Group: "handler"})
	if err == nil {
		for _, v := range keywords {
			if aliasMap[v.Roomid] == nil {
				aliasMap[v.Roomid] = map[string]string{}
			}
			aliasMap[v.Roomid][v.Target] = v.Phrase
		}
	}

	// 生成指令菜单
	helper := []string{}
	for _, v := range handlers {
		if v.Level > 0 {
			if up == nil || v.Level > up.Level {
				continue // 没有权限
			}
		}
		if (msg.IsGroup && v.RoomAble) || (!msg.IsGroup && v.ChatAble) {
			cmd := v.Command
			if aliasMap[msg.Roomid] != nil && aliasMap[msg.Roomid][v.Command] != "" {
				cmd = aliasMap[msg.Roomid][v.Command]
			} else if aliasMap["-"] != nil && aliasMap["-"][v.Command] != "" {
				cmd = aliasMap["-"][v.Command]
			}
			helper = append(helper, fmt.Sprintf("【%s】%s", cmd, v.Describe))
		}
	}

	// 数组转为字符串
	text := strings.Join(helper, "\n") + "\n"

	// 自定义帮助信息
	if len(setting.HelpAdditive) > 1 {
		text += setting.HelpAdditive + "\n"
	}

	// 分割线
	text += "----------------\n"

	// 当前用户状态信息
	if up.Level > 0 {
		text += fmt.Sprintf("级别 %d；", up.Level)
	}

	// 对话模型相关配置
	llmCount, _ := llmodel.Count(&llmodel.CountParam{})
	if llmCount > 0 {
		model := aichat.UserModel(msg.Sender, msg.Roomid).Family
		if len(model) > 1 {
			text += fmt.Sprintf("对话模型 %s；", model)
		}
		text += fmt.Sprintf("上下文长度 %d/%d；", aichat.CountHistory(msg.Sender, msg.Roomid), setting.ModelHistory)
	}

	return text + "祝你好运！"

}
